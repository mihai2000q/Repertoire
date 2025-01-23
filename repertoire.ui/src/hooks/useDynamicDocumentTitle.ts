import { useAppDispatch, useAppSelector } from '../state/store.ts'
import { setDocumentTitle } from '../state/globalSlice.ts'

export default function useDynamicDocumentTitle() {
  const dispatch = useAppDispatch()
  const title = useAppSelector((state) => state.global.documentTitle)
  return (newTitle: string | ((prevTitle: string) => string)) =>
    dispatch(setDocumentTitle(typeof newTitle === 'string' ? newTitle : newTitle(title)))
}
