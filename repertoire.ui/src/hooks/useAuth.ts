export default function useAuth(): boolean {
  const token = localStorage.getItem('token')
  return token !== null
}
