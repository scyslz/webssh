import request from '@/utils/request'
export function checkSSH(sshInfo) {
    return request.get(`/check?sshInfo=${sshInfo}`)
}

export function getSSHList() {
    return request.get("/ssh/list");
}

export function saveSSHList(sshList) {
    return request.post("/ssh/save", sshList);
}
