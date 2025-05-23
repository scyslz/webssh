import { getSSHList } from '@/api/common'
import state from './state';

export default {
    setLanguage({ commit }, language) {
        commit('SET_LANGUAGE', language)
    },
    async fetchSshList({ commit }) {
        try {
            const sshList = await getSSHList();
            if (!sshList || sshList.length === 0) {
                return;
            }
            const encodedSshList = window.btoa(JSON.stringify(sshList)); // Assuming getSSHList returns an array of objects
            console.log(encodedSshList)
            commit('SET_LIST', encodedSshList);
            var cache=  localStorage.getItem("sshList")
           console.log(JSON.parse( window.atob(cache)))

        } catch (error) {
            console.error('Error fetching SSH list:', error);
            // Handle the error appropriately
        }
    }
}
