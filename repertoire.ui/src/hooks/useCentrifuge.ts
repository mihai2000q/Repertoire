import { Centrifuge } from 'centrifuge'
import { useEffect } from 'react'
import { useLazyGetCentrifugeTokenQuery } from '../state/api/authApi.ts'

let centrifuge: Centrifuge | undefined

export default function useCentrifuge(): Centrifuge | undefined {
  const [getNewToken] = useLazyGetCentrifugeTokenQuery()

  useEffect(() => {
    if (!centrifuge) {
      centrifuge = new Centrifuge(import.meta.env.VITE_CENTRIFUGO_URL, {
        getToken: async () => await getNewToken().unwrap()
      })
    }
    centrifuge.connect()
    return () => centrifuge.disconnect()
  }, [])

  return centrifuge
}
