import { defineStore } from 'pinia'
import { ref } from 'vue'
import { client } from '../utils/client'

export const useSideMenuStore = defineStore('side-menu', () => {
    const CHECKING_TEXT = 'Checking...'
    const health = ref(CHECKING_TEXT)
    const condensed = ref(false)

    const checkHealth = async () => {
        health.value = CHECKING_TEXT
        health.value = await client.api.v1HealthList()
            .then(res => res.data.status)
            .catch(_ => 'Can\'t reach')
            .finally(() => setTimeout(() => {
                checkHealth()
            }, 20000))
    }

    return {
        condensed,
        health,
        checkHealth,
    }
})
