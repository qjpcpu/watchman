package router

const SYS_ID = "_alfred_"

// default policy
func defaultPolicyFunc(from_id, body, to_id string) bool {
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
func defaultSocketFunc() string {
    return "/tmp/router.socket"
}
func DefaultBuilder() Builder {
    return Builder{
        defaultPolicyFunc,
        defaultSocketFunc,
    }
}
