import { useContext } from 'react'
import { MainScrollContext, MainScrollContextReturnType } from '../context/MainScrollContext.tsx'

export default function useMainScroll() {
  return useContext<MainScrollContextReturnType>(MainScrollContext)
}
