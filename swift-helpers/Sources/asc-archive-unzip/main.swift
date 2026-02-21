import ArgumentParser
import Compression
import Foundation

// MARK: - Errors

enum ArchiveUnzipError: Error, LocalizedError {
    case invalidArchive(String)
    case extractionFailed(String)
    case fileNotFound(String)
    
    var errorDescription: String? {
        switch self {
        case .invalidArchive(let detail):
            return "Invalid archive: \(detail)"
        case .extractionFailed(let detail):
            return "Extraction failed: \(detail)"
        case .fileNotFound(let path):
            return "File not found: \(path)"
        }
    }
}

// MARK: - Extraction Result

struct ExtractionResult: Codable {
    let success: Bool
    let archivePath: String
    let outputPath: String
    let filesExtracted: Int
    let totalSize: Int64
    let duration: Double
}

// MARK: - Archive Extraction

func extractArchive(
    archivePath: String,
    outputPath: String,
    progress: Bool = false
) throws -> ExtractionResult {
    let fm = FileManager.default
    let startTime = Date()
    
    // Ensure archive exists
    guard fm.fileExists(atPath: archivePath) else {
        throw ArchiveUnzipError.fileNotFound(archivePath)
    }
    
    // Create output directory
    try fm.createDirectory(atPath: outputPath, withIntermediateDirectories: true)
    
    // Use unzip command with progress
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/unzip")
    
    var args = ["-o"]  // Overwrite existing files
    
    if !progress {
        args.append("-q")  // Quiet mode
    }
    
    args.append(contentsOf: [archivePath, "-d", outputPath])
    
    task.arguments = args
    
    let pipe = Pipe()
    task.standardOutput = pipe
    task.standardError = pipe
    
    try task.run()
    task.waitUntilExit()
    
    if task.terminationStatus != 0 {
        let errorData = pipe.fileHandleForReading.readDataToEndOfFile()
        let errorMessage = String(data: errorData, encoding: .utf8) ?? "Unknown error"
        throw ArchiveUnzipError.extractionFailed(errorMessage)
    }
    
    // Count extracted files and size
    let (fileCount, totalSize) = try calculateExtractedStats(at: outputPath)
    
    let duration = Date().timeIntervalSince(startTime)
    
    return ExtractionResult(
        success: true,
        archivePath: archivePath,
        outputPath: outputPath,
        filesExtracted: fileCount,
        totalSize: totalSize,
        duration: duration
    )
}

func calculateExtractedStats(at path: String) throws -> (Int, Int64) {
    let fm = FileManager.default
    var fileCount = 0
    var totalSize: Int64 = 0
    
    let enumerator = fm.enumerator(atPath: path)
    while let file = enumerator?.nextObject() as? String {
        let filePath = (path as NSString).appendingPathComponent(file)
        if let attributes = try? fm.attributesOfItem(atPath: filePath),
           let type = attributes[.type] as? FileAttributeType,
           type == .typeRegular {
            fileCount += 1
            if let size = attributes[.size] as? Int64 {
                totalSize += size
            }
        }
    }
    
    return (fileCount, totalSize)
}

func listArchiveContents(archivePath: String) throws -> [String] {
    guard FileManager.default.fileExists(atPath: archivePath) else {
        throw ArchiveUnzipError.fileNotFound(archivePath)
    }
    
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/unzip")
    task.arguments = ["-l", archivePath]
    
    let pipe = Pipe()
    task.standardOutput = pipe
    try task.run()
    task.waitUntilExit()
    
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    let output = String(data: data, encoding: .utf8) ?? ""
    
    var files: [String] = []
    let lines = output.components(separatedBy: .newlines)
    
    for line in lines {
        // Parse lines like: "  1234  01-01-2024 12:00   filename.ext"
        let trimmed = line.trimmingCharacters(in: .whitespaces)
        let components = trimmed.components(separatedBy: .whitespaces).filter { !$0.isEmpty }
        
        if components.count >= 4,
           Int(components[0]) != nil,  // Size
           components[1].contains("-") {  // Date
            let filename = components.dropFirst(3).joined(separator: " ")
            if !filename.isEmpty && !filename.hasSuffix("/") {
                files.append(filename)
            }
        }
    }
    
    return files
}

// MARK: - Commands

struct ExtractCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "extract",
        abstract: "Extract a .zip or .ipa archive"
    )
    
    @Option(name: .long, help: "Archive path")
    var archive: String
    
    @Option(name: .long, help: "Output directory")
    var output: String
    
    @Flag(name: .long, help: "Show progress")
    var progress: Bool = false
    
    mutating func run() throws {
        let result = try extractArchive(
            archivePath: archive,
            outputPath: output,
            progress: progress
        )
        
        let encoder = JSONEncoder()
        encoder.outputFormatting = .sortedKeys
        let data = try encoder.encode(result)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct ListCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "list",
        abstract: "List archive contents without extracting"
    )
    
    @Argument(help: "Archive path")
    var archive: String
    
    mutating func run() throws {
        let files = try listArchiveContents(archivePath: archive)
        
        let dict: [String: Any] = [
            "archive": archive,
            "files": files,
            "count": files.count
        ]
        
        let data = try JSONSerialization.data(withJSONObject: dict, options: .sortedKeys)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct InfoCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "info",
        abstract: "Get archive information"
    )
    
    @Argument(help: "Archive path")
    var archive: String
    
    mutating func run() throws {
        let fm = FileManager.default
        
        guard fm.fileExists(atPath: archive) else {
            throw ArchiveUnzipError.fileNotFound(archive)
        }
        
        let size = (try? fm.attributesOfItem(atPath: archive)[.size] as? Int64) ?? 0
        let files = try? listArchiveContents(archivePath: archive)
        
        let dict: [String: Any] = [
            "path": archive,
            "size_bytes": size,
            "size_mb": Double(size) / 1024.0 / 1024.0,
            "file_count": files?.count ?? 0
        ]
        
        let data = try JSONSerialization.data(withJSONObject: dict, options: .sortedKeys)
        print(String(data: data, encoding: .utf8)!)
    }
}

@main
struct ArchiveUnzipCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "asc-archive-unzip",
        abstract: "Fast archive extraction for .xcarchive, .ipa, and .zip files",
        version: "0.1.0",
        subcommands: [ExtractCommand.self, ListCommand.self, InfoCommand.self]
    )
}
