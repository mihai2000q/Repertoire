import { useEffect } from 'react'
import { useAppDispatch, useAppSelector } from '../state/store.ts'
import { api } from '../state/api.ts'
import useCentrifuge from './useCentrifuge.ts'

export default function useSearchQueryCacheInvalidation() {
  const dispatch = useAppDispatch()
  const centrifuge = useCentrifuge()

  const userID = useAppSelector((state) => state.global.userID)

  useEffect(() => {
    if (!userID) return () => {}

    const channel = `search:${userID}`
    const sub = centrifuge.getSubscription(channel) ?? centrifuge.newSubscription(channel)

    sub.on('publication', (data) => {
      if (data.data.action === 'SEARCH_CACHE_INVALIDATION') {
        dispatch(api.util.invalidateTags(['Search']))
      }
    })

    sub.subscribe()
    return () => sub.unsubscribe()
  }, [dispatch, userID])
}
