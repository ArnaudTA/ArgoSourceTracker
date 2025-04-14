import { Api } from "../api/Api";

export const client = new Api({
    baseURL: `${import.meta.env.VITE_API_URL}`
})