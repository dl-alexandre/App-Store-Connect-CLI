import ArgumentParser
import Foundation
import Security

// MARK: - Errors

enum KeychainHelperError: Error, LocalizedError, Equatable {
    case accessDenied(OSStatus)
    case interactionNotAllowed(OSStatus)
    case itemNotFound
    case duplicateItem
    case encodingFailed(String)
    case decodingFailed(String)
    case unexpectedStatus(OSStatus)
    case invalidInput(String)
    
    var errorDescription: String? {
        switch self {
        case .accessDenied(let code):
            return "Keychain access denied (errSecAuthFailed \(code))"
        case .interactionNotAllowed(let code):
            return "Keychain interaction not allowed (errSecInteractionNotAllowed \(code))"
        case .itemNotFound:
            return "Keychain item not found"
        case .duplicateItem:
            return "Keychain item already exists"
        case .encodingFailed(let detail):
            return "Failed to encode: \(detail)"
        case .decodingFailed(let detail):
            return "Failed to decode: \(detail)"
        case .unexpectedStatus(let code):
            return "Unexpected keychain status: \(code)"
        case .invalidInput(let detail):
            return "Invalid input: \(detail)"
        }
    }
}

func makeKeychainError(_ status: OSStatus) -> KeychainHelperError {
    switch status {
    case -25293:
        return .accessDenied(status)
    case -25308:
        return .interactionNotAllowed(status)
    case -25300:
        return .itemNotFound
    case -25299:
        return .duplicateItem
    default:
        return .unexpectedStatus(status)
    }
}

// MARK: - Constants

let keychainServiceName = "asc"
let keychainItemPrefix = "asc:credential:"

// MARK: - Credential Structure

struct CredentialPayload: Codable {
    let keyID: String
    let issuerID: String
    let privateKeyPath: String
    
    enum CodingKeys: String, CodingKey {
        case keyID = "key_id"
        case issuerID = "issuer_id"
        case privateKeyPath = "private_key_path"
    }
}

// MARK: - Keychain Operations

func storeCredential(name: String, keyID: String, issuerID: String, privateKeyPath: String) throws {
    let payload = CredentialPayload(keyID: keyID, issuerID: issuerID, privateKeyPath: privateKeyPath)
    let data = try JSONEncoder().encode(payload)
    
    let account = keychainItemPrefix + name
    let label = "ASC API Key (\(name))"
    
    let query: [CFString: Any] = [
        kSecClass: kSecClassGenericPassword,
        kSecAttrService: keychainServiceName,
        kSecAttrAccount: account,
        kSecAttrLabel: label,
        kSecValueData: data,
        kSecAttrAccessible: kSecAttrAccessibleWhenUnlocked,
        kSecAttrSynchronizable: kCFBooleanFalse as Any
    ]
    
    let status = SecItemAdd(query as CFDictionary, nil)
    
    if status == -25299 { // errSecDuplicateItem
        // Update existing item
        let updateQuery: [CFString: Any] = [
            kSecClass: kSecClassGenericPassword,
            kSecAttrService: keychainServiceName,
            kSecAttrAccount: account
        ]
        let attributes: [CFString: Any] = [kSecValueData: data]
        let updateStatus = SecItemUpdate(updateQuery as CFDictionary, attributes as CFDictionary)
        guard updateStatus == errSecSuccess else {
            throw makeKeychainError(updateStatus)
        }
    } else if status != errSecSuccess {
        throw makeKeychainError(status)
    }
}

func getCredential(name: String) throws -> CredentialPayload? {
    let account = keychainItemPrefix + name
    
    let query: [CFString: Any] = [
        kSecClass: kSecClassGenericPassword,
        kSecAttrService: keychainServiceName,
        kSecAttrAccount: account,
        kSecReturnData: kCFBooleanTrue as Any,
        kSecMatchLimit: kSecMatchLimitOne
    ]
    
    var result: AnyObject?
    let status = SecItemCopyMatching(query as CFDictionary, &result)
    
    if status == -25300 {
        return nil
    }
    
    guard status == errSecSuccess else {
        throw makeKeychainError(status)
    }
    
    guard let data = result as? Data else {
        throw KeychainHelperError.decodingFailed("Unexpected result type")
    }
    
    return try JSONDecoder().decode(CredentialPayload.self, from: data)
}

func listCredentials() throws -> [(name: String, payload: CredentialPayload)] {
    let query: [CFString: Any] = [
        kSecClass: kSecClassGenericPassword,
        kSecAttrService: keychainServiceName,
        kSecReturnAttributes: kCFBooleanTrue as Any,
        kSecReturnData: kCFBooleanTrue as Any,
        kSecMatchLimit: kSecMatchLimitAll
    ]
    
    var result: AnyObject?
    let status = SecItemCopyMatching(query as CFDictionary, &result)
    
    if status == -25300 {
        return []
    }
    
    guard status == errSecSuccess else {
        throw makeKeychainError(status)
    }
    
    guard let items = result as? [[String: Any]] else {
        return []
    }
    
    var credentials: [(name: String, payload: CredentialPayload)] = []
    
    for item in items {
        guard let account = item[kSecAttrAccount as String] as? String,
              account.hasPrefix(keychainItemPrefix),
              let data = item[kSecValueData as String] as? Data else {
            continue
        }
        
        let name = String(account.dropFirst(keychainItemPrefix.count))
        
        if let payload = try? JSONDecoder().decode(CredentialPayload.self, from: data) {
            credentials.append((name: name, payload: payload))
        }
    }
    
    return credentials
}

func deleteCredential(name: String) throws {
    let account = keychainItemPrefix + name
    
    let query: [CFString: Any] = [
        kSecClass: kSecClassGenericPassword,
        kSecAttrService: keychainServiceName,
        kSecAttrAccount: account
    ]
    
    let status = SecItemDelete(query as CFDictionary)
    
    guard status == errSecSuccess || status == -25300 else {
        throw makeKeychainError(status)
    }
}

// MARK: - Commands

struct StoreCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "store",
        abstract: "Store credentials in the keychain"
    )
    
    @Argument(help: "Credential name")
    var name: String
    
    @Option(name: .long, help: "Key ID")
    var keyID: String
    
    @Option(name: .long, help: "Issuer ID")
    var issuerID: String
    
    @Option(name: .long, help: "Path to private key file")
    var privateKeyPath: String
    
    func run() throws {
        try storeCredential(name: name, keyID: keyID, issuerID: issuerID, privateKeyPath: privateKeyPath)
        print("Stored credential '\(name)'")
    }
}

struct GetCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "get",
        abstract: "Retrieve credentials from the keychain"
    )
    
    @Argument(help: "Credential name")
    var name: String
    
    @Option(name: .long, help: "Output format: json, yaml")
    var format: String = "json"
    
    func run() throws {
        guard let payload = try getCredential(name: name) else {
            throw KeychainHelperError.itemNotFound
        }
        
        let dict: [String: String] = [
            "name": name,
            "key_id": payload.keyID,
            "issuer_id": payload.issuerID,
            "private_key_path": payload.privateKeyPath
        ]
        
        let data = try JSONSerialization.data(withJSONObject: dict, options: .sortedKeys)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct ListCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "list",
        abstract: "List all stored credentials"
    )
    
    @Option(name: .long, help: "Output format: json, names")
    var format: String = "names"
    
    func run() throws {
        let credentials = try listCredentials()
        
        switch format {
        case "json":
            let array = credentials.map { cred -> [String: String] in
                [
                    "name": cred.name,
                    "key_id": cred.payload.keyID,
                    "issuer_id": cred.payload.issuerID
                ]
            }
            let data = try JSONSerialization.data(withJSONObject: array, options: .sortedKeys)
            print(String(data: data, encoding: .utf8)!)
        default:
            for cred in credentials {
                print(cred.name)
            }
        }
    }
}

struct DeleteCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "delete",
        abstract: "Delete credentials from the keychain"
    )
    
    @Argument(help: "Credential name")
    var name: String
    
    @Flag(name: .long, help: "Skip confirmation")
    var force: Bool = false
    
    func run() throws {
        try deleteCredential(name: name)
        print("Deleted credential '\(name)'")
    }
}

@main
struct KeychainCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "asc-keychain",
        abstract: "Native macOS keychain operations for ASC credentials using Security.framework",
        version: "0.1.0",
        subcommands: [StoreCommand.self, GetCommand.self, ListCommand.self, DeleteCommand.self]
    )
}
