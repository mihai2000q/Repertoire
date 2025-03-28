import { useLazyGetCentrifugeTokenQuery } from '../state/api.ts'
import { Centrifuge } from 'centrifuge'
import { useEffect } from 'react'

let centrifuge: Centrifuge | undefined

export default function useCentrifuge(): Centrifuge {
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

  return centrifuge! // Non-null assertion (we ensure it exists)
}
