// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "asc-helpers",
    platforms: [
        .macOS(.v14)
    ],
    products: [
        .executable(
            name: "asc-jwt-sign",
            targets: ["asc-jwt-sign"]
        ),
        .executable(
            name: "asc-keychain",
            targets: ["asc-keychain"]
        ),
        .executable(
            name: "asc-screenshot-frame",
            targets: ["asc-screenshot-frame"]
        ),
        .executable(
            name: "asc-image-optimize",
            targets: ["asc-image-optimize"]
        ),
        .executable(
            name: "asc-bundle-validate",
            targets: ["asc-bundle-validate"]
        ),
        .executable(
            name: "asc-ipa-pack",
            targets: ["asc-ipa-pack"]
        ),
        .executable(
            name: "asc-simulator",
            targets: ["asc-simulator"]
        ),
        .executable(
            name: "asc-video-encode",
            targets: ["asc-video-encode"]
        ),
        .executable(
            name: "asc-codesign",
            targets: ["asc-codesign"]
        ),
        .executable(
            name: "asc-archive-unzip",
            targets: ["asc-archive-unzip"]
        )
    ],
    dependencies: [
        .package(url: "https://github.com/apple/swift-argument-parser.git", from: "1.3.0")
    ],
    targets: [
        // JWT signing helper - CryptoKit hardware-accelerated
        .executableTarget(
            name: "asc-jwt-sign",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-jwt-sign"
        ),
        
        // Keychain helper - Security.framework native access
        .executableTarget(
            name: "asc-keychain",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-keychain"
        ),
        
        // Screenshot framing helper - Core Image/Metal accelerated
        .executableTarget(
            name: "asc-screenshot-frame",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-screenshot-frame"
        ),
        
        // Image optimization helper - Core Image/Metal accelerated
        .executableTarget(
            name: "asc-image-optimize",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-image-optimize"
        ),
        
        // Bundle validation helper - CodeSigning + Security.framework
        .executableTarget(
            name: "asc-bundle-validate",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-bundle-validate"
        ),
        
        // IPA packaging helper - libcompression
        .executableTarget(
            name: "asc-ipa-pack",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-ipa-pack"
        ),
        
        // Simulator helper - XCTest + CoreSimulator
        .executableTarget(
            name: "asc-simulator",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-simulator"
        ),
        
        // Video encoding helper - AVFoundation
        .executableTarget(
            name: "asc-video-encode",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-video-encode"
        ),
        
        // Code signing helper - Security.framework
        .executableTarget(
            name: "asc-codesign",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-codesign"
        ),
        
        // Archive extraction helper - libcompression
        .executableTarget(
            name: "asc-archive-unzip",
            dependencies: [
                .product(name: "ArgumentParser", package: "swift-argument-parser")
            ],
            path: "Sources/asc-archive-unzip"
        ),
        
        // Test targets
        .testTarget(
            name: "JWTHelperTests",
            dependencies: ["asc-jwt-sign"],
            path: "Tests/JWTHelperTests"
        ),
        .testTarget(
            name: "KeychainHelperTests",
            dependencies: ["asc-keychain"],
            path: "Tests/KeychainHelperTests"
        ),
        .testTarget(
            name: "ScreenshotFrameTests",
            dependencies: ["asc-screenshot-frame"],
            path: "Tests/ScreenshotFrameTests"
        )
    ]
)
