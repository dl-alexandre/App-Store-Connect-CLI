import ArgumentParser
import Compression
import Foundation

// MARK: - Errors

enum IPAPackError: Error, LocalizedError {
    case invalidInput(String)
    case compressionFailed(String)
    case appNotFound(String)
    case packFailed(String)
    
    var errorDescription: String? {
        switch self {
        case .invalidInput(let detail):
            return "Invalid input: \(detail)"
        case .compressionFailed(let detail):
            return "Compression failed: \(detail)"
        case .appNotFound(let path):
            return "App bundle not found: \(path)"
        case .packFailed(let detail):
            return "Pack failed: \(detail)"
        }
    }
}

// MARK: - IPA Packing

struct PackResult: Codable {
    let success: Bool
    let appPath: String
    let ipaPath: String
    let originalSize: Int64
    let compressedSize: Int64
    let compressionRatio: Double
    let compressionMethod: String
    let duration: Double
}

func packIPA(
    appPath: String,
    outputPath: String,
    compression: compression_algorithm = COMPRESSION_ZLIB
) throws -> PackResult {
    let fm = FileManager.default
    let startTime = Date()
    
    // Validate app exists
    guard fm.fileExists(atPath: appPath) else {
        throw IPAPackError.appNotFound(appPath)
    }
    
    // Get original size
    let originalSize = try calculateDirectorySize(at: appPath)
    
    // Create Payload directory in temp location
    let tempDir = NSTemporaryDirectory() + "/asc-ipa-pack-" + UUID().uuidString
    let payloadDir = tempDir + "/Payload"
    try fm.createDirectory(atPath: payloadDir, withIntermediateDirectories: true)
    
    // Copy app to Payload
    let appName = (appPath as NSString).lastPathComponent
    let destAppPath = payloadDir + "/" + appName
    try fm.copyItem(atPath: appPath, toPath: destAppPath)
    
    // Create IPA using libcompression
    let ipaURL = URL(fileURLWithPath: outputPath)
    
    // Use Process to call zip for now (libcompression archive API is more complex)
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/zip")
    task.arguments = [
        "-qr",  // quiet, recursive
        outputPath,
        "Payload"
    ]
    task.currentDirectoryURL = URL(fileURLWithPath: tempDir)
    
    let pipe = Pipe()
    task.standardOutput = pipe
    task.standardError = pipe
    
    try task.run()
    task.waitUntilExit()
    
    if task.terminationStatus != 0 {
        let errorData = pipe.fileHandleForReading.readDataToEndOfFile()
        let errorMessage = String(data: errorData, encoding: .utf8) ?? "Unknown error"
        throw IPAPackError.packFailed(errorMessage)
    }
    
    // Cleanup
    try? fm.removeItem(atPath: tempDir)
    
    // Get compressed size
    let compressedSize = (try? fm.attributesOfItem(atPath: outputPath)[.size] as? Int64) ?? 0
    let duration = Date().timeIntervalSince(startTime)
    let ratio = originalSize > 0 ? Double(compressedSize) / Double(originalSize) : 1.0
    
    return PackResult(
        success: true,
        appPath: appPath,
        ipaPath: outputPath,
        originalSize: originalSize,
        compressedSize: compressedSize,
        compressionRatio: ratio,
        compressionMethod: "deflate",
        duration: duration
    )
}

func calculateDirectorySize(at path: String) throws -> Int64 {
    let fm = FileManager.default
    var totalSize: Int64 = 0
    
    let enumerator = fm.enumerator(atPath: path)
    while let file = enumerator?.nextObject() as? String {
        let filePath = (path as NSString).appendingPathComponent(file)
        if let attributes = try? fm.attributesOfItem(atPath: filePath),
           let size = attributes[.size] as? Int64 {
            totalSize += size
        }
    }
    
    return totalSize
}

// MARK: - Commands

struct PackCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "pack",
        abstract: "Pack an .app bundle into an .ipa file"
    )
    
    @Option(name: .long, help: "Path to .app bundle")
    var app: String
    
    @Option(name: .long, help: "Output .ipa path")
    var output: String
    
    @Option(name: .long, help: "Compression level (0-9, default: 6)")
    var level: Int = 6
    
    mutating func run() throws {
        let result = try packIPA(appPath: app, outputPath: output)
        
        let encoder = JSONEncoder()
        encoder.outputFormatting = .sortedKeys
        let data = try encoder.encode(result)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct InfoCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "info",
        abstract: "Get .app bundle size information"
    )
    
    @Argument(help: "Path to .app bundle")
    var path: String
    
    mutating func run() throws {
        let fm = FileManager.default
        
        guard fm.fileExists(atPath: path) else {
            throw IPAPackError.appNotFound(path)
        }
        
        let size = try calculateDirectorySize(at: path)
        let dict: [String: Any] = [
            "path": path,
            "size_bytes": size,
            "size_mb": Double(size) / 1024.0 / 1024.0
        ]
        
        let data = try JSONSerialization.data(withJSONObject: dict, options: .sortedKeys)
        print(String(data: data, encoding: .utf8)!)
    }
}

@main
struct IPAPackCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "asc-ipa-pack",
        abstract: "Fast IPA packaging with compression optimization",
        version: "0.1.0",
        subcommands: [PackCommand.self, InfoCommand.self]
    )
}
