import ArgumentParser
import CryptoKit
import Foundation

// MARK: - Errors

enum JWTSignError: Error, LocalizedError {
    case invalidPrivateKey(String)
    case keyFileReadError(String)
    case invalidIssuerID
    case invalidKeyID
    case signingFailed(String)
    
    var errorDescription: String? {
        switch self {
        case .invalidPrivateKey(let detail):
            return "Invalid private key: \(detail)"
        case .keyFileReadError(let detail):
            return "Failed to read key file: \(detail)"
        case .invalidIssuerID:
            return "Invalid or missing issuer ID"
        case .invalidKeyID:
            return "Invalid or missing key ID"
        case .signingFailed(let detail):
            return "Signing failed: \(detail)"
        }
    }
}

// MARK: - JWT Structures

struct JWTHeader: Encodable {
    let alg: String = "ES256"
    let kid: String
    let typ: String = "JWT"
}

struct JWTClaims: Encodable {
    let iss: String
    let iat: Int
    let exp: Int
    let aud: String
}

// MARK: - Helper Functions

/// Token lifetime: 10 minutes (matches Go implementation)
let jwtTokenLifetime: TimeInterval = 10 * 60

func base64URLEncode(_ data: Data) -> String {
    data.base64EncodedString()
        .replacingOccurrences(of: "+", with: "-")
        .replacingOccurrences(of: "/", with: "_")
        .replacingOccurrences(of: "=", with: "")
}

func loadPrivateKey(from path: String) throws -> P256.Signing.PrivateKey {
    let url = URL(fileURLWithPath: path)
    let pemData: Data
    do {
        pemData = try Data(contentsOf: url)
    } catch {
        throw JWTSignError.keyFileReadError(error.localizedDescription)
    }
    
    guard let pemString = String(data: pemData, encoding: .utf8) else {
        throw JWTSignError.keyFileReadError("File is not valid UTF-8")
    }
    
    // Extract base64 content from PEM
    let lines = pemString.components(separatedBy: .newlines)
    let base64Lines = lines.filter { !$0.hasPrefix("-") && !$0.isEmpty }
    let base64String = base64Lines.joined()
    
    guard let keyData = Data(base64Encoded: base64String) else {
        throw JWTSignError.invalidPrivateKey("Failed to decode base64 content")
    }
    
    do {
        // Try PKCS8 first (common for .p8 files)
        return try P256.Signing.PrivateKey(x963Representation: keyData)
    } catch {
        // Fall back to SEC1 format
        do {
            return try P256.Signing.PrivateKey(rawRepresentation: keyData)
        } catch {
            throw JWTSignError.invalidPrivateKey("Key is not valid P-256 format: \(error)")
        }
    }
}

func generateJWT(issuerID: String, keyID: String, privateKey: P256.Signing.PrivateKey) throws -> String {
    let now = Date()
    let iat = Int(now.timeIntervalSince1970)
    let exp = Int(now.addingTimeInterval(jwtTokenLifetime).timeIntervalSince1970)
    
    let header = JWTHeader(kid: keyID)
    let claims = JWTClaims(
        iss: issuerID,
        iat: iat,
        exp: exp,
        aud: "appstoreconnect-v1"
    )
    
    let encoder = JSONEncoder()
    encoder.outputFormatting = .sortedKeys
    
    let headerData = try encoder.encode(header)
    let payloadData = try encoder.encode(claims)
    
    let headerEncoded = base64URLEncode(headerData)
    let payloadEncoded = base64URLEncode(payloadData)
    let signingInput = "\(headerEncoded).\(payloadEncoded)"
    
    guard let data = signingInput.data(using: .utf8) else {
        throw JWTSignError.signingFailed("Failed to encode signing input")
    }
    
    let signature = try privateKey.signature(for: data)
    let signatureEncoded = base64URLEncode(signature.rawRepresentation)
    
    return "\(signingInput).\(signatureEncoded)"
}

/// Signs a string payload using ES256 (ECDSA P-256 + SHA-256).
func sign(payload: String, with privateKey: P256.Signing.PrivateKey) throws -> String {
    guard let data = payload.data(using: .utf8) else {
        throw JWTSignError.signingFailed("Could not encode signing input as UTF-8")
    }
    do {
        let signature = try privateKey.signature(for: data)
        return base64URLEncode(signature.rawRepresentation)
    } catch {
        throw JWTSignError.signingFailed(error.localizedDescription)
    }
}

// MARK: - Command

@main
struct JWTSignCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "asc-jwt-sign",
        abstract: "Generate JWT tokens for App Store Connect API authentication using CryptoKit hardware acceleration",
        version: "0.1.0"
    )
    
    @Option(name: .long, help: "App Store Connect Issuer ID")
    var issuerID: String
    
    @Option(name: .long, help: "App Store Connect Key ID")
    var keyID: String
    
    @Option(name: .long, help: "Path to private key file (.p8)")
    var privateKeyPath: String
    
    @Option(name: .long, help: "Output format: token (default), json")
    var output: String = "token"
    
    @Flag(name: .long, help: "Validate the generated token without output")
    var validate: Bool = false
    
    mutating func run() throws {
        // Validate inputs
        guard !issuerID.isEmpty else {
            throw JWTSignError.invalidIssuerID
        }
        guard !keyID.isEmpty else {
            throw JWTSignError.invalidKeyID
        }
        
        // Load private key
        let privateKey = try loadPrivateKey(from: privateKeyPath)
        
        // Generate JWT
        let token = try generateJWT(issuerID: issuerID, keyID: keyID, privateKey: privateKey)
        
        // Output result
        switch output.lowercased() {
        case "json":
            let result: [String: Any] = [
                "token": token,
                "expires_in": Int(jwtTokenLifetime)
            ]
            let jsonData = try JSONSerialization.data(withJSONObject: result, options: .sortedKeys)
            print(String(data: jsonData, encoding: .utf8)!)
        default:
            print(token)
        }
    }
}
