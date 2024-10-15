import { useAppSelector } from '../state/store'

export default function useAuth(): boolean {
  const token = useAppSelector((state) => state.auth.token)
  return token !== null
}
