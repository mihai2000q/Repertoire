import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../state/store'
import { useEffect } from 'react'
import { setErrorPath } from '../state/slice/globalSlice.ts'

export default function useErrorRedirection(): void {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()

  const errorPath = useAppSelector((state) => state.global.errorPath)
  useEffect(() => {
    if (errorPath) {
      navigate(errorPath)
      dispatch(setErrorPath(undefined))
    }
  }, [dispatch, errorPath, navigate])
}
