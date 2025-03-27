import { useEffect } from 'react'
import { useAppDispatch, useAppSelector } from '../state/store.ts'
import { api } from '../state/api.ts'
import useCentrifuge from './useCentrifuge.ts'

export default function useSearchQueryCacheInvalidation() {
  const dispatch = useAppDispatch()

  const userID = useAppSelector((state) => state.global.userID)

  const [centrifuge, isLoading] = useCentrifuge()

  useEffect(() => {
    if (!userID || isLoading) return

    const sub = centrifuge.newSubscription(`search:${userID}`)

    sub.on('publication', (data) => {
      if (data.data.action === 'SEARCH_CACHE_INVALIDATION') {
        dispatch(api.util.invalidateTags(['Search']))
      }
    })

    sub.subscribe()
    centrifuge.connect()
    return () => {
      centrifuge.disconnect()
      sub.unsubscribe()
    }
  }, [dispatch, userID, isLoading])
}
