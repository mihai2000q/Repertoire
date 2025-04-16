import { useShallowEffect } from '@mantine/hooks'
import { DependencyList, EffectCallback, useEffect, useRef } from 'react'

export default function useDidUpdateShallow(fn: EffectCallback, dependencies?: DependencyList) {
  const mounted = useRef(false)

  useEffect(
    () => () => {
      mounted.current = false
    },
    []
  )

  useShallowEffect(() => {
    if (mounted.current) return fn()
    mounted.current = true
    return undefined
  }, dependencies)
}
