package utils




func IsValidShortCode(code string) bool {
    if len(code) < 6 || len(code) > 12 {
        return false
    }
    
    // Base64 URL encoding uses: A-Z, a-z, 0-9, - (minus), and _ (underscore)
    for _, c := range code {
        if !(c >= 'a' && c <= 'z') && 
           !(c >= 'A' && c <= 'Z') && 
           !(c >= '0' && c <= '9') &&
           c != '-' &&  // Allow minus
           c != '_' {   // Allow underscore
            return false
        }
    }
    return true
}