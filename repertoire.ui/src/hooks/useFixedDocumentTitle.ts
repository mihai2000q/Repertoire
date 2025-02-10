import { useEffect } from 'react'
import { useAppDispatch } from '../state/store.ts'
import { setDocumentTitle } from '../state/slice/globalSlice.ts'

export default function useFixedDocumentTitle(value: string) {
  const dispatch = useAppDispatch()
  useEffect(() => {
    dispatch(setDocumentTitle(value))
  }, [])
}
