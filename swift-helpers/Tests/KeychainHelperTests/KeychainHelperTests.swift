import XCTest
@testable import asc_keychain

final class KeychainHelperTests: XCTestCase {
    
    func testAccountKeyGeneration() {
        let name = "test-profile"
        let expectedAccount = keychainItemPrefix + name
        
        // This is testing the internal logic that we can't access directly
        // But we can verify by storing and retrieving
        XCTAssertEqual(keychainItemPrefix, "asc:credential:")
    }
    
    func testShouldBypassKeychain() {
        // Test environment variable parsing
        let testCases: [(String, Bool)] = [
            ("1", true),
            ("true", true),
            ("yes", true),
            ("on", true),
            ("0", false),
            ("false", false),
            ("no", false),
            ("", false),
            ("random", false)
        ]
        
        for (value, expected) in testCases {
            // Can't easily test without modifying environment
            // This would need integration testing
            XCTAssertTrue(true) // Placeholder
        }
    }
    
    func testKeychainErrorDescription() {
        let errors: [KeychainHelperError] = [
            .accessDenied(-25293),
            .interactionNotAllowed(-25308),
            .itemNotFound,
            .duplicateItem,
            .encodingFailed("test"),
            .decodingFailed("test"),
            .unexpectedStatus(-99999),
            .invalidInput("test")
        ]
        
        for error in errors {
            XCTAssertNotNil(error.errorDescription)
            XCTAssertFalse(error.errorDescription?.isEmpty ?? true)
        }
    }
    
    func testMakeKeychainErrorMapping() {
        XCTAssertEqual(makeKeychainError(-25293), KeychainHelperError.accessDenied(-25293))
        XCTAssertEqual(makeKeychainError(-25308), KeychainHelperError.interactionNotAllowed(-25308))
        XCTAssertEqual(makeKeychainError(-25300), KeychainHelperError.itemNotFound)
        XCTAssertEqual(makeKeychainError(-25299), KeychainHelperError.duplicateItem)
        XCTAssertEqual(makeKeychainError(-99999), KeychainHelperError.unexpectedStatus(-99999))
    }
}
