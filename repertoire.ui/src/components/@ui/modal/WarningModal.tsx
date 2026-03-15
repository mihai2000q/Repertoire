import { Button, Group, Modal, ModalProps, Stack, Text } from '@mantine/core'
import { ReactNode } from 'react'

interface WarningModalProps extends ModalProps {
  opened: boolean
  onClose: () => void
  title: string
  description: string | ReactNode
  onYes: () => void
  isLoading?: boolean
}

function WarningModal({
  opened,
  onClose,
  title,
  description,
  onYes,
  isLoading,
  ...props
}: WarningModalProps) {
  function internalOnYes() {
    onYes()
    onClose()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={title} centered {...props}>
      <Stack px={'xs'} py={0}>
        {typeof description === 'string' ? <Text fw={500}>{description}</Text> : description}
        <Group gap={'xxs'} style={{ alignSelf: 'end' }}>
          <Button variant={'subtle'} onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={internalOnYes} loading={isLoading}>
            Yes
          </Button>
        </Group>
      </Stack>
    </Modal>
  )
}

export default WarningModal
