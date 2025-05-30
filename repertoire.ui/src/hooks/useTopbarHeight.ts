import useAuth from './useAuth.ts'

export default function useTopbarHeight() {
  return useAuth() ? '65px' : '0px'
}
