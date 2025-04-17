import { ref } from 'vue'
import { defineStore } from 'pinia'
import { ApplicationApplicationStatus, ApplicationApplicationSummary } from '../api/Api'
import { client } from '../utils/client'

export const useHomeStore = defineStore('home-menu', () => {
    const filters = ref<ApplicationApplicationStatus[]>([])
    const applications = ref<ApplicationApplicationSummary[]>([])
    const isLoading = ref(true)
    const errorMessage = ref('')

    const fetchApps = async () => {
        isLoading.value = true
        errorMessage.value = ''
        client.api.v1AppsList({
            filter: filters.value.join(',')
        })
            .then(res => applications.value = res.data)
            .catch(reason => errorMessage.value = reason)
            .finally(() => isLoading.value = false)
    }
    return {
        filters,
        applications,
        isLoading,
        errorMessage,
        fetchApps
    }
})