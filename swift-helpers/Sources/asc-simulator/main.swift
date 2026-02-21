import ArgumentParser
import Foundation

// MARK: - Errors

enum SimulatorError: Error, LocalizedError {
    case deviceNotFound(String)
    case screenshotFailed(String)
    case installFailed(String)
    case launchFailed(String)
    case listDevicesFailed(String)
    
    var errorDescription: String? {
        switch self {
        case .deviceNotFound(let id):
            return "Device not found: \(id)"
        case .screenshotFailed(let detail):
            return "Screenshot failed: \(detail)"
        case .installFailed(let detail):
            return "Install failed: \(detail)"
        case .launchFailed(let detail):
            return "Launch failed: \(detail)"
        case .listDevicesFailed(let detail):
            return "Failed to list devices: \(detail)"
        }
    }
}

// MARK: - Simulator Types

struct SimulatorDevice: Codable {
    let udid: String
    let name: String
    let deviceType: String
    let runtime: String
    let state: String
    let isAvailable: Bool
}

struct ScreenshotResult: Codable {
    let success: Bool
    let deviceUDID: String
    let outputPath: String
    let timestamp: String
    let size: Int64
}

// MARK: - Simulator Control

func listSimulators() throws -> [SimulatorDevice] {
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/xcrun")
    task.arguments = ["simctl", "list", "-j", "devices"]
    
    let pipe = Pipe()
    task.standardOutput = pipe
    try task.run()
    task.waitUntilExit()
    
    let data = pipe.fileHandleForReading.readDataToEndOfFile()
    
    guard let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
          let devicesByRuntime = json["devices"] as? [String: [[String: String]]] else {
        throw SimulatorError.listDevicesFailed("Could not parse simctl output")
    }
    
    var devices: [SimulatorDevice] = []
    
    for (runtime, runtimeDevices) in devicesByRuntime {
        for deviceInfo in runtimeDevices {
            if let udid = deviceInfo["udid"],
               let name = deviceInfo["name"],
               let state = deviceInfo["state"],
               let availability = deviceInfo["isAvailable"] {
                devices.append(SimulatorDevice(
                    udid: udid,
                    name: name,
                    deviceType: deviceInfo["deviceTypeIdentifier"] ?? "unknown",
                    runtime: runtime,
                    state: state,
                    isAvailable: availability == "true"
                ))
            }
        }
    }
    
    return devices
}

func findSimulator(name: String? = nil, udid: String? = nil) throws -> SimulatorDevice {
    let devices = try listSimulators()
    
    if let udid = udid {
        guard let device = devices.first(where: { $0.udid == udid }) else {
            throw SimulatorError.deviceNotFound(udid)
        }
        return device
    }
    
    if let name = name {
        // Try exact match first
        if let device = devices.first(where: { $0.name == name && $0.isAvailable }) {
            return device
        }
        
        // Try partial match
        if let device = devices.first(where: { 
            $0.name.contains(name) && $0.isAvailable 
        }) {
            return device
        }
        
        throw SimulatorError.deviceNotFound(name)
    }
    
    // Return first available booted device, or first available
    if let booted = devices.first(where: { $0.state == "Booted" && $0.isAvailable }) {
        return booted
    }
    
    if let available = devices.first(where: { $0.isAvailable }) {
        return available
    }
    
    throw SimulatorError.deviceNotFound("No available simulator found")
}

func installApp(deviceUDID: String, appPath: String) throws {
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/xcrun")
    task.arguments = ["simctl", "install", deviceUDID, appPath]
    
    let pipe = Pipe()
    task.standardError = pipe
    try task.run()
    task.waitUntilExit()
    
    if task.terminationStatus != 0 {
        let errorData = pipe.fileHandleForReading.readDataToEndOfFile()
        let errorMessage = String(data: errorData, encoding: .utf8) ?? "Unknown error"
        throw SimulatorError.installFailed(errorMessage)
    }
}

func launchApp(deviceUDID: String, bundleID: String) throws {
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/xcrun")
    task.arguments = ["simctl", "launch", deviceUDID, bundleID]
    
    try task.run()
    task.waitUntilExit()
    
    if task.terminationStatus != 0 {
        throw SimulatorError.launchFailed("App launch failed with exit code \(task.terminationStatus)")
    }
}

func takeScreenshot(deviceUDID: String, outputPath: String) throws -> ScreenshotResult {
    let task = Process()
    task.executableURL = URL(fileURLWithPath: "/usr/bin/xcrun")
    task.arguments = ["simctl", "io", deviceUDID, "screenshot", outputPath]
    
    let pipe = Pipe()
    task.standardError = pipe
    try task.run()
    task.waitUntilExit()
    
    if task.terminationStatus != 0 {
        let errorData = pipe.fileHandleForReading.readDataToEndOfFile()
        let errorMessage = String(data: errorData, encoding: .utf8) ?? "Unknown error"
        throw SimulatorError.screenshotFailed(errorMessage)
    }
    
    // Get file size
    let fm = FileManager.default
    let size = (try? fm.attributesOfItem(atPath: outputPath)[.size] as? Int64) ?? 0
    
    let formatter = ISO8601DateFormatter()
    
    return ScreenshotResult(
        success: true,
        deviceUDID: deviceUDID,
        outputPath: outputPath,
        timestamp: formatter.string(from: Date()),
        size: size
    )
}

// MARK: - Commands

struct ListCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "list",
        abstract: "List available simulators"
    )
    
    @Flag(name: .long, help: "Show only booted devices")
    var booted: Bool = false
    
    mutating func run() throws {
        let devices = try listSimulators()
        let filtered = booted ? devices.filter { $0.state == "Booted" } : devices
        
        let encoder = JSONEncoder()
        encoder.outputFormatting = .sortedKeys
        let data = try encoder.encode(filtered)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct ScreenshotCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "screenshot",
        abstract: "Take a screenshot of a simulator"
    )
    
    @Option(name: .long, help: "Device UDID (or uses first booted)")
    var udid: String?
    
    @Option(name: .long, help: "Device name (e.g., 'iPhone 16 Pro')")
    var device: String?
    
    @Option(name: .long, help: "Output path")
    var output: String
    
    mutating func run() throws {
        let sim = try findSimulator(name: device, udid: udid)
        let result = try takeScreenshot(deviceUDID: sim.udid, outputPath: output)
        
        let encoder = JSONEncoder()
        encoder.outputFormatting = .sortedKeys
        let data = try encoder.encode(result)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct InstallCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "install",
        abstract: "Install an app on a simulator"
    )
    
    @Option(name: .long, help: "Device UDID or name")
    var udid: String?
    
    @Option(name: .long, help: "Device name")
    var device: String?
    
    @Argument(help: "Path to .app bundle")
    var appPath: String
    
    mutating func run() throws {
        let sim = try findSimulator(name: device, udid: udid)
        try installApp(deviceUDID: sim.udid, appPath: appPath)
        
        let dict: [String: String] = [
            "status": "success",
            "device": sim.name,
            "udid": sim.udid,
            "app": appPath
        ]
        
        let data = try JSONSerialization.data(withJSONObject: dict, options: .sortedKeys)
        print(String(data: data, encoding: .utf8)!)
    }
}

struct LaunchCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "launch",
        abstract: "Launch an app on a simulator"
    )
    
    @Option(name: .long, help: "Device UDID or name")
    var udid: String?
    
    @Option(name: .long, help: "Device name")
    var device: String?
    
    @Argument(help: "Bundle identifier")
    var bundleID: String
    
    mutating func run() throws {
        let sim = try findSimulator(name: device, udid: udid)
        try launchApp(deviceUDID: sim.udid, bundleID: bundleID)
        
        let dict: [String: String] = [
            "status": "success",
            "device": sim.name,
            "udid": sim.udid,
            "bundle_id": bundleID
        ]
        
        let data = try JSONSerialization.data(withJSONObject: dict, options: .sortedKeys)
        print(String(data: data, encoding: .utf8)!)
    }
}

@main
struct SimulatorCommand: ParsableCommand {
    static let configuration = CommandConfiguration(
        commandName: "asc-simulator",
        abstract: "iOS Simulator control for screenshot capture and app testing",
        version: "0.1.0",
        subcommands: [ListCommand.self, ScreenshotCommand.self, InstallCommand.self, LaunchCommand.self]
    )
}
