package router

const SYS_ID = "_alfred_"

// default policy
func DefaultPolicy(from_id, body, to_id string) bool {
    if from_id == to_id {
        return false
    }
    if from_id == SYS_ID {
        return true
    }
    if to_id == SYS_ID {
        return true
    } else {
        return false
    }
}
