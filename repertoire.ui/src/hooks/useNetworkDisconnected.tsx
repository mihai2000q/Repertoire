import { useNetwork } from '@mantine/hooks'
import { useEffect, useRef } from 'react'
import { toast } from 'react-toastify'
import { Group, Loader, Text } from '@mantine/core'

export default function useNetworkDisconnected() {
  const { online } = useNetwork()
  const toastId = useRef(null)
  useEffect(() => {
    if (!online) {
      toastId.current = toast.error(
        () => (
          <Group gap={'xs'}>
            <Text fw={500}>Please check your internet connection</Text>
            <Loader type={'dots'} size={28} color={'gray'} />
          </Group>
        ),
        { autoClose: false, closeButton: false }
      )
    } else if (toastId.current !== null) {
      toast.dismiss(toastId.current)
      toast.success('The internet is back! If the page is not loading, please refresh!')
      toastId.current = null
    }
  }, [online])
}
