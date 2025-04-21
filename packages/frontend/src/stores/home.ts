import { computed, ComputedRef, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { TypesApplicationStatus, TypesSummary } from '../api/Api'
import { client } from '../utils/client'
import { PageState } from 'primevue'

export const useHomeStore = defineStore('home-menu', () => {
    const StatusChoice = ref<{ name: TypesApplicationStatus, key: TypesApplicationStatus }[]>([
        { name: TypesApplicationStatus.Ignored, key: TypesApplicationStatus.Ignored },
        { name: TypesApplicationStatus.Outdated, key: TypesApplicationStatus.Outdated },
        { name: TypesApplicationStatus.UpToDate, key: TypesApplicationStatus.UpToDate },
    ]);
    const selectedStatus = ref<{ name: TypesApplicationStatus, key: TypesApplicationStatus }[]>([
        { name: TypesApplicationStatus.UpToDate, key: TypesApplicationStatus.UpToDate },
        { name: TypesApplicationStatus.Outdated, key: TypesApplicationStatus.Outdated },
    ]);

    const isLoading = ref(true)
    const errorMessage = ref('')
    const total = ref(0)
    const limit = ref(10)
    const page = ref(1)
    const offset = ref(0)

    const applications = ref<TypesSummary[]>([])

    function refreshApps() {
        isLoading.value = true
        errorMessage.value = ''
        client.api.v1AppsList({
            filter: selectedStatus.value
                .map(status => status.key)
                .join(','),
            offset: offset.value,
            limit: limit.value,
        })
            .then(res => {
                applications.value = res.data
                total.value = Number(res.headers["x-total"])
            })
            .catch(reason => errorMessage.value = reason)
            .finally(() => isLoading.value = false)
    }

    function setPage(p: PageState) {
        console.log(p);

        offset.value = p.first
        limit.value = p.rows
        page.value = p.page
        refreshApps()
    }

    type Meter = {
        label: TypesApplicationStatus, color: string, value: number
    }
    const numberOfStatus: Record<TypesApplicationStatus, ComputedRef<number>> = {
        [TypesApplicationStatus.UpToDate]: computed(() => applications.value.filter(app => app.status === TypesApplicationStatus.UpToDate).length),
        [TypesApplicationStatus.Ignored]: computed(() => applications.value.filter(app => app.status === TypesApplicationStatus.Ignored).length),
        [TypesApplicationStatus.Outdated]: computed(() => applications.value.filter(app => app.status === TypesApplicationStatus.Outdated).length),
    }
    const statusColor: Record<TypesApplicationStatus, string> = {
        Outdated: '#ffb223',
        "Up-to-date": '#02e702',
        Ignored: '#d3d3d38f'
    }
    const meters = computed<Meter[]>(() => Object.entries(numberOfStatus)
        .filter(([_k, v]) => v.value > 0)
        .map(([k, v]) => ({ label: k as TypesApplicationStatus, value: v.value, color: statusColor[k as TypesApplicationStatus] }))
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
    }
})