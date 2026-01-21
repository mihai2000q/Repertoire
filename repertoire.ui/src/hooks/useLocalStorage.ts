import { useState } from 'react'
import { useDidUpdate } from '@mantine/hooks'

interface useLocalStorageOptions<T> {
  key: string
  defaultValue?: T
  serialize?: (value: T) => string
  deserialize?: (item: string) => T
}

export default function useLocalStorage<T>({
  key,
  defaultValue,
  serialize,
  deserialize
}: useLocalStorageOptions<T>): [T, (value: T) => void] {
  const [item, setItem] = useState<T>(() => {
    const item = localStorage.getItem(key)
    if (item === null || item === undefined) return defaultValue
    return deserialize ? deserialize(item) : JSON.parse(item)
  })

  useDidUpdate(() => {
    if (item !== undefined)
      localStorage.setItem(key, serialize ? serialize(item) : JSON.stringify(item))
    else localStorage.removeItem(key)
  }, [item])

  return [item, setItem]
}
