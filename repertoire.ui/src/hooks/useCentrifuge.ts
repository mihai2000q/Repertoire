import { useGetCentrifugeTokenQuery, useLazyGetCentrifugeTokenQuery } from '../state/api.ts'
import { Centrifuge } from 'centrifuge'

export default function useCentrifuge(): [Centrifuge, boolean] {
  const { data: token, isLoading } = useGetCentrifugeTokenQuery()

  async function refreshToken(): Promise<string> {
    const [trigger, { data: newToken }] = useLazyGetCentrifugeTokenQuery()
    await trigger().unwrap()
    return newToken
  }

  if (isLoading)
    return [undefined, isLoading]

  return [
    new Centrifuge(import.meta.env.VITE_CENTRIFUGO_URL, {
      token: token,
      getToken: refreshToken // refreshes token when it expires
    }),
    isLoading
  ]
}
