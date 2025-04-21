import type { PageState } from 'primevue'
import type { ComputedRef } from 'vue'
import type { TypesSummary } from '../api/Api'
import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'
import { TypesApplicationStatus } from '../api/Api'
import { client } from '../utils/client'

export const useHomeStore = defineStore('home-menu', () => {
    const StatusChoice = ref<{ name: TypesApplicationStatus, key: TypesApplicationStatus }[]>(Object.values(TypesApplicationStatus).map(v => ({
        name: v as TypesApplicationStatus,
        key: v as TypesApplicationStatus,
    })))
    const selectedStatus = ref<TypesApplicationStatus[]>([])

    const isLoading = ref(true)
    const errorMessage = ref('')
    const limit = ref(25)
    const page = ref(1)
    const offset = ref(0)
    const stats = ref<Record<string, number>>({})
    const total = computed(() => Object.values(stats.value).reduce((sum, n) => (sum + n), 0))
    const applications = ref<TypesSummary[]>([])

    function refreshApps() {
        isLoading.value = true
        errorMessage.value = ''

        client.api.v1AppsList({
            filter: selectedStatus.value.join(','),
            offset: offset.value,
            limit: limit.value,
        })
            .then((res) => {
                applications.value = res.data.items
                stats.value = res.data.stats
            })
            .catch(reason => errorMessage.value = reason)
            .finally(() => isLoading.value = false)
    }

    function setPage(p: PageState) {
        offset.value = p.first
        limit.value = p.rows
        page.value = p.page
        refreshApps()
    }

    interface Meter {
        label: TypesApplicationStatus
        color: string
        value: number
    }
    const numberOfStatus: Record<TypesApplicationStatus, ComputedRef<number>> = {
        [TypesApplicationStatus.UpToDate]: computed(() => applications.value.filter(app => app.status === TypesApplicationStatus.UpToDate).length),
        [TypesApplicationStatus.Ignored]: computed(() => applications.value.filter(app => app.status === TypesApplicationStatus.Ignored).length),
        [TypesApplicationStatus.Outdated]: computed(() => applications.value.filter(app => app.status === TypesApplicationStatus.Outdated).length),
        [TypesApplicationStatus.Error]: computed(() => applications.value.filter(app => app.status === TypesApplicationStatus.Error).length),
    }
    const statusColor: Record<TypesApplicationStatus, string> = {
        'Up-to-date': '#02e702',
        'Outdated': '#ffb223',
        'Ignored': '#d3d3d38f',
        'Error': '#ff0000',
    }
    const meters = computed<Meter[]>(() => Object.entries(numberOfStatus)
        .filter(([_k, v]) => v.value > 0)
        .map(([k, v]) => ({ label: k as TypesApplicationStatus, value: v.value, color: statusColor[k as TypesApplicationStatus] })),
    )
    watch(selectedStatus, refreshApps)
    return {
        StatusChoice,
        applications,
        isLoading,
        errorMessage,
        page,
        selectedStatus,
        total,
        limit,
        setPage,
        refreshApps,
        meters,
        stats,
    }
}, {
    persist: {
        storage: sessionStorage,
        pick: ['selectedStatus', 'limit'],
    },
})
