export default function useIsDesktop() {
    return import.meta.env.VITE_PLATFORM === 'desktop'
}