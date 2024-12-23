export default function useIsProduction() {
  return import.meta.env.VITE_ENVIRONMENT === 'production'
}
