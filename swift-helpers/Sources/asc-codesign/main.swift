import ArgumentParser
import Foundation
import Security

// MARK: - Errors

enum CodeSignError: Error, LocalizedError {
    case invalidInput(String)
    case signingFailed(String)
    case verificationFailed(String)
    case certificateNotFound(String)
    
    var errorDescription: String? {
        switch self {
        case .invalidInput(let detail):
            return "Invalid input: \(detail)"
        case .signingFailed(let detail):
            return "Signing failed: \(detail)"
        case .verificationFailed(let detail):
            return "Verification failed: \(detail)"
        case .certificateNotFound(let name):
            return "Certificate not found: \(name)"
        }
    }
}

// MARK: - Signing Result

struct SigningResult: Codable {
    let success: Bool
    let path: String
    let identity: String?
    let timestamp: String?
    let errors: [String]
}

struct VerificationResult: Codable {
    let valid: Bool
    let path: String
    let authority: String?
    let identifier: String?
    let teamIdentifier: String?
    let timestamp: String?
}

// MARK: - Code Signing

func signApp(
    path: String,
    identity: String? = nil,
    entitlements: String? = nil,
    force: Bool = false
) throws -> SigningResult {
    var errors: [String] = []
    
    var args = ["--sign"]
    
    if let identity = identity {
        args.append(identity)
    } else {
        args.append("-") // Ad-hoc signing
    }
    
    if force {
        args.append("--force")
    }
    
    if let entitlements = entitlements {
        args.append(contentsOf: ["--entitlements", entitlements])
    }
    
    args.append("--timestamp")
    args.append(path)
    
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/codesign")
    task.arguments = args
    
    let errorPipe = Pipe()
    task.standardError = errorPipe
    
    try task.run()
    task.waitUntilExit()
    
    if task.terminationStatus != 0 {
        let errorData = errorPipe.fileHandleForReading.readDataToEndOfFile()
        let errorMessage = String(data: errorData, encoding: .utf8) ?? "Unknown error"
        errors.append(errorMessage)
        throw CodeSignError.signingFailed(errorMessage)
    }
    
    // Get signature info
    let verifyResult = try? verifySignature(path: path)
    
    return SigningResult(
        success: true,
        path: path,
        identity: identity,
        timestamp: verifyResult?.timestamp,
        errors: errors
    )
}

func verifySignature(path: String) throws -> VerificationResult {
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/codesign")
    task.arguments = ["-dvv", path]
    
    let errorPipe = Pipe()
    task.standardError = errorPipe
    
    try task.run()
    task.waitUntilExit()
    
    let errorData = errorPipe.fileHandleForReading.readDataToEndOfFile()
    let output = String(data: errorData, encoding: .utf8) ?? ""
    
    var authority: String?
    var identifier: String?
    var teamIdentifier: String?
    var timestamp: String?
    
    let lines = output.components(separatedBy: .newlines)
    for line in lines {
        if line.contains("Authority=") {
            authority = line.components(separatedBy: "=").last?.trimmingCharacters(in: .whitespaces)
        } else if line.contains("Identifier=") {
            identifier = line.components(separatedBy: "=").last?.trimmingCharacters(in: .whitespaces)
        } else if line.contains("TeamIdentifier=") {
            teamIdentifier = line.components(separatedBy: "=").last?.trimmingCharacters(in: .whitespaces)
        } else if line.contains("Timestamp=") {
            timestamp = line.components(separatedBy: "=").last?.trimmingCharacters(in: .whitespaces)
        }
    }
    
    let valid = task.terminationStatus == 0 && !output.contains("error")
    
    return VerificationResult(
        valid: valid,
        path: path,
        authority: authority,
        identifier: identifier,
        teamIdentifier: teamIdentifier,
        timestamp: timestamp
    )
}

func listIdentities() throws -> [[String: String]] {
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/security")
    task.arguments = ["find-identity", "-v", "-p", "codesigning"]
    
    let pipe = Pipe()
    task.standardOutput = pipe
    try task.run()
    task.waitUntilExit()
    
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    let output = String(data: data, encoding: .utf8) ?? ""
    
    var identities: [[String: String]] = []
    let lines = output.components(separatedBy: .newlines)
    
    for line in lines {
        // Parse lines like: "1) ABCDEF123... "iPhone Developer: John Doe (TEAMID)"
        if line.contains(")") {
            let components = line.components(separatedBy: ")")
            if components.count >= 2 {
                let hash = components[1].trimmingCharacters(in: .whitespaces).prefix(40)
                let nameStart = line.range(of: "\"")
                let nameEnd = line.range(of: "\"", options: .backwards)
                
                if let start = nameStart?.upperBound, let end = nameEnd?.lowerBound {
                    let name = String(line[start..<end])
                    identities.append([
                        "hash": String(hash),
                        "name": name
                    ])
                }
            }
        }
    }
    
    return identities
}

// MARK: - Commands

struct SignCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "sign",
        abstract: "Sign an app bundle"
    )
    
    @Argument(help: "Path to .app bundle")
    var path: String
    
    @Option(name: .long, help: "Signing identity (hash or name)")
    var identity: String?
    
    @Option(name: .long, help: "Path to entitlements file")
    var entitlements: String?
    
    @Flag(name: .long, help: "Force re-signing")
    var force: Bool = false
    
    mutating func run() throws {
        let result = try signApp(
            path: path,
            identity: identity,
            entitlements: entitlements,
            force: force
        )
        
        let encoder = JSONEncoder()
        encoder.outputFormatting = .sortedKeys
        let data = try encoder.encode(result)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct VerifyCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "verify",
        abstract: "Verify code signature"
    )
    
    @Argument(help: "Path to signed bundle")
    var path: String
    
    mutating func run() throws {
        let result = try verifySignature(path: path)
        
        let encoder = JSONEncoder()
        encoder.outputFormatting = .sortedKeys
        let data = try encoder.encode(result)
        print(String(data: data, encoding: .utf8)!)
        
        if !result.valid {
            throw ExitCode(1)
        }
    }
}

struct ListCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "list",
        abstract: "List available signing identities"
    )
    
    mutating func run() throws {
        let identities = try listIdentities()
        
        let data = try JSONSerialization.data(withJSONObject: identities, options: .sortedKeys)
        print(String(data: data, encoding: .utf8)!)
    }
}

@main
struct CodeSignCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "asc-codesign",
        abstract: "Code signing utilities for iOS/macOS apps",
        version: "0.1.0",
        subcommands: [SignCommand.self, VerifyCommand.self, ListCommand.self]
    )
}
