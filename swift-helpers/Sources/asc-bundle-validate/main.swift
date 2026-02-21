import ArgumentParser
import Foundation
import Security

// MARK: - Errors

enum BundleValidateError: Error, LocalizedError {
    case invalidBundle(String)
    case validationFailed(String)
    case signatureInvalid(String)
    case provisioningProfileExpired(String)
    case entitlementsMismatch(String)
    
    var errorDescription: String? {
        switch self {
        case .invalidBundle(let detail):
            return "Invalid bundle: \(detail)"
        case .validationFailed(let detail):
            return "Validation failed: \(detail)"
        case .signatureInvalid(let detail):
            return "Invalid signature: \(detail)"
        case .provisioningProfileExpired(let date):
            return "Provisioning profile expired on \(date)"
        case .entitlementsMismatch(let detail):
            return "Entitlements mismatch: \(detail)"
        }
    }
}

// MARK: - Validation Result

struct ValidationIssue: Codable {
    let severity: String  // "error", "warning", "info"
    let code: String
    let message: String
    let path: String?
}

struct ValidationResult: Codable {
    let valid: Bool
    let bundlePath: String
    let bundleIdentifier: String?
    let bundleVersion: String?
    let issues: [ValidationIssue]
    let signature: SignatureInfo?
    let provisioningProfile: ProvisioningProfileInfo?
    let entitlements: [String: AnyCodable]?
}

struct SignatureInfo: Codable {
    let valid: Bool
    let authority: String?
    let identifier: String?
    let timestamp: String?
}

struct ProvisioningProfileInfo: Codable {
    let valid: Bool
    let name: String?
    let appID: String?
    let expirationDate: String?
    let teamID: String?
    let isExpired: Bool
}

// Helper to encode Any values
struct AnyCodable: Codable {
    let value: Any
    
    init(_ value: Any) {
        self.value = value
    }
    
    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        if let string = value as? String {
            try container.encode(string)
        } else if let int = value as? Int {
            try container.encode(int)
        } else if let bool = value as? Bool {
            try container.encode(bool)
        } else if let array = value as? [String] {
            try container.encode(array)
        } else {
            try container.encode(String(describing: value))
        }
    }
    
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if let string = try? container.decode(String.self) {
            value = string
        } else if let int = try? container.decode(Int.self) {
            value = int
        } else if let bool = try? container.decode(Bool.self) {
            value = bool
        } else if let array = try? container.decode([String].self) {
            value = array
        } else {
            value = ""
        }
    }
}

// MARK: - Validation Logic

func validateBundle(at path: String) throws -> ValidationResult {
    let fm = FileManager.default
    var issues: [ValidationIssue] = []
    
    // Check if path exists
    guard fm.fileExists(atPath: path) else {
        throw BundleValidateError.invalidBundle("Path does not exist: \(path)")
    }
    
    // Determine bundle type
    let isApp = path.hasSuffix(".app")
    let isIPA = path.hasSuffix(".ipa")
    
    var bundlePath = path
    var cleanupPath: String?
    
    // Extract IPA if needed
    if isIPA {
        let tempDir = NSTemporaryDirectory() + "/asc-bundle-validate-" + UUID().uuidString
        try fm.createDirectory(atPath: tempDir, withIntermediateDirectories: true)
        cleanupPath = tempDir
        
        // Extract Payload directory
        let task = Process()
        task.executableURL = URL(fileURLWithPath: "/usr/bin/unzip")
        task.arguments = ["-q", path, "-d", tempDir]
        try task.run()
        task.waitUntilExit()
        
        // Find .app in Payload
        let payloadPath = tempDir + "/Payload"
        let contents = try fm.contentsOfDirectory(atPath: payloadPath)
        if let appName = contents.first(where: { $0.hasSuffix(".app") }) {
            bundlePath = payloadPath + "/" + appName
        } else {
            issues.append(ValidationIssue(
                severity: "error",
                code: "NO_APP_IN_IPA",
                message: "No .app bundle found in IPA",
                path: nil
            ))
        }
    }
    
    // Read Info.plist
    let infoPlistPath = bundlePath + "/Info.plist"
    var bundleIdentifier: String?
    var bundleVersion: String?
    var entitlements: [String: AnyCodable]?
    
    if fm.fileExists(atPath: infoPlistPath) {
        if let plistData = fm.contents(atPath: infoPlistPath),
           let plist = try? PropertyListSerialization.propertyList(from: plistData, format: nil) as? [String: Any] {
            bundleIdentifier = plist["CFBundleIdentifier"] as? String
            bundleVersion = plist["CFBundleShortVersionString"] as? String
            
            if bundleIdentifier == nil {
                issues.append(ValidationIssue(
                    severity: "error",
                    code: "MISSING_BUNDLE_ID",
                    message: "CFBundleIdentifier not found in Info.plist",
                    path: infoPlistPath
                ))
            }
        } else {
            issues.append(ValidationIssue(
                severity: "error",
                code: "INVALID_INFO_PLIST",
                message: "Could not parse Info.plist",
                path: infoPlistPath
            ))
        }
    } else {
        issues.append(ValidationIssue(
            severity: "error",
            code: "MISSING_INFO_PLIST",
            message: "Info.plist not found",
            path: nil
        ))
    }
    
    // Validate signature using codesign
    var signatureInfo: SignatureInfo?
    let codesignTask = Process()
    codesignTask.executableURL = URL(fileURLWithPath: "/usr/bin/codesign")
    codesignTask.arguments = ["-dvv", bundlePath]
    
    let outputPipe = Pipe()
    let errorPipe = Pipe()
    codesignTask.standardOutput = outputPipe
    codesignTask.standardError = errorPipe
    
    try codesignTask.run()
    codesignTask.waitUntilExit()
    
    let errorData = errorPipe.fileHandleForReading.readDataToEndOfFile()
    let errorOutput = String(data: errorData, encoding: .utf8) ?? ""
    
    signatureInfo = parseSignatureOutput(errorOutput)
    
    if codesignTask.terminationStatus != 0 {
        issues.append(ValidationIssue(
            severity: "error",
            code: "SIGNATURE_INVALID",
            message: "Code signature validation failed",
            path: bundlePath
        ))
    }
    
    // Check embedded provisioning profile
    var profileInfo: ProvisioningProfileInfo?
    let profilePath = bundlePath + "/embedded.mobileprovision"
    
    if fm.fileExists(atPath: profilePath) {
        profileInfo = try validateProvisioningProfile(at: profilePath)
        if profileInfo?.isExpired ?? false {
            issues.append(ValidationIssue(
                severity: "error",
                code: "PROFILE_EXPIRED",
                message: "Provisioning profile has expired",
                path: profilePath
            ))
        }
    } else {
        issues.append(ValidationIssue(
            severity: "warning",
            code: "NO_PROVISIONING_PROFILE",
            message: "No embedded provisioning profile found (may be App Store build)",
            path: nil
        ))
    }
    
    // Check entitlements
    let entitlementsTask = Process()
    entitlementsTask.executableURL = URL(fileURLWithPath: "/usr/bin/codesign")
    entitlementsTask.arguments = ["--entitlements", "-", "--xml", "-d", bundlePath]
    
    let entitlementsPipe = Pipe()
    entitlementsTask.standardOutput = entitlementsPipe
    try entitlementsTask.run()
    entitlementsTask.waitUntilExit()
    
    let entitlementsData = entitlementsPipe.fileHandleForReading.readDataToEndOfFile()
    if let entitlementsPlist = try? PropertyListSerialization.propertyList(from: entitlementsData, format: nil) as? [String: Any] {
        entitlements = entitlementsPlist.mapValues { AnyCodable($0) }
    }
    
    // Cleanup temp directory if created
    if let cleanup = cleanupPath {
        try? fm.removeItem(atPath: cleanup)
    }
    
    return ValidationResult(
        valid: issues.filter { $0.severity == "error" }.isEmpty,
        bundlePath: path,
        bundleIdentifier: bundleIdentifier,
        bundleVersion: bundleVersion,
        issues: issues,
        signature: signatureInfo,
        provisioningProfile: profileInfo,
        entitlements: entitlements
    )
}

func parseSignatureOutput(_ output: String) -> SignatureInfo {
    var authority: String?
    var identifier: String?
    var timestamp: String?
    
    let lines = output.components(separatedBy: .newlines)
    for line in lines {
        if line.contains("Authority=") {
            authority = line.components(separatedBy: "=").last?.trimmingCharacters(in: .whitespaces)
        } else if line.contains("Identifier=") {
            identifier = line.components(separatedBy: "=").last?.trimmingCharacters(in: .whitespaces)
        } else if line.contains("Timestamp=") {
            timestamp = line.components(separatedBy: "=").last?.trimmingCharacters(in: .whitespaces)
        }
    }
    
    return SignatureInfo(
        valid: !output.contains("error"),
        authority: authority,
        identifier: identifier,
        timestamp: timestamp
    )
}

func validateProvisioningProfile(at path: String) throws -> ProvisioningProfileInfo {
    // Read the provisioning profile (CMS signed plist)
    let data = try Data(contentsOf: URL(fileURLWithPath: path))
    
    // Use security command to extract plist
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/security")
    task.arguments = ["cms", "-D", "-i", path]
    
    let pipe = Pipe()
    task.standardOutput = pipe
    try task.run()
    task.waitUntilExit()
    
    let plistData = pipe.fileHandleForReading.readDataToEndOfFile()
    
    guard let plist = try? PropertyListSerialization.propertyList(from: plistData, format: nil) as? [String: Any] else {
        throw BundleValidateError.validationFailed("Could not parse provisioning profile")
    }
    
    let name = plist["Name"] as? String
    let teamID = (plist["TeamIdentifier"] as? [String])?.first
    let expirationString = plist["ExpirationDate"] as? String
    let appID = plist["AppIDName"] as? String
    
    let isExpired: Bool
    if let expString = expirationString,
       let expDate = parseProvisioningDate(expString) {
        isExpired = expDate < Date()
    } else {
        isExpired = false
    }
    
    return ProvisioningProfileInfo(
        valid: !isExpired,
        name: name,
        appID: appID,
        expirationDate: expirationString,
        teamID: teamID,
        isExpired: isExpired
    )
}

func parseProvisioningDate(_ string: String) -> Date? {
    // Try multiple date formats
    let formatters = [
        "yyyy-MM-dd'T'HH:mm:ss'Z'",
        "yyyy-MM-dd HH:mm:ss Z",
        "EEE MMM dd HH:mm:ss z yyyy"
    ]
    
    for format in formatters {
        let formatter = DateFormatter()
        formatter.dateFormat = format
        formatter.locale = Locale(identifier: "en_US_POSIX")
        if let date = formatter.date(from: string) {
            return date
        }
    }
    
    return nil
}

// MARK: - Commands

struct ValidateCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "validate",
        abstract: "Validate an app bundle or IPA"
    )
    
    @Argument(help: "Path to .app bundle or .ipa file")
    var path: String
    
    @Flag(name: .long, help: "Fail on warnings")
    var strict: Bool = false
    
    mutating func run() throws {
        let result = try validateBundle(at: path)
        
        let encoder = JSONEncoder()
        encoder.outputFormatting = .sortedKeys
        let data = try encoder.encode(result)
        print(String(data: data, encoding: .utf8)!)
        
        // Exit with error code if validation failed or strict mode with warnings
        if !result.valid || (strict && !result.issues.isEmpty) {
            throw ExitCode(1)
        }
    }
}

struct InfoCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "info",
        abstract: "Get bundle information without validation"
    )
    
    @Argument(help: "Path to .app bundle")
    var path: String
    
    mutating func run() throws {
        let fm = FileManager.default
        let infoPlistPath = path + "/Info.plist"
        
        guard fm.fileExists(atPath: infoPlistPath),
              let data = fm.contents(atPath: infoPlistPath),
              let plist = try? PropertyListSerialization.propertyList(from: data, format: nil) as? [String: Any] else {
            throw BundleValidateError.invalidBundle("Could not read Info.plist")
        }
        
        let dict: [String: Any] = [
            "bundle_id": plist["CFBundleIdentifier"] as? String ?? "",
            "version": plist["CFBundleShortVersionString"] as? String ?? "",
            "build": plist["CFBundleVersion"] as? String ?? "",
            "name": plist["CFBundleDisplayName"] as? String ?? plist["CFBundleName"] as? String ?? "",
            "minimum_os": plist["MinimumOSVersion"] as? String ?? "",
            "platform": plist["DTPlatformName"] as? String ?? ""
        ]
        
        let jsonData = try JSONSerialization.data(withJSONObject: dict, options: .sortedKeys)
        print(String(data: jsonData, encoding: .utf8)!)
    }
}

@main
struct BundleValidateCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "asc-bundle-validate",
        abstract: "Validate iOS/macOS app bundles and IPAs for App Store submission",
        version: "0.1.0",
        subcommands: [ValidateCommand.self, InfoCommand.self]
    )
}
