import { toast } from 'vue3-toastify';

export function useNotify(message) {
    toast(message, {
        autoClose: 1000,
        theme: 'dark',
    });
}
