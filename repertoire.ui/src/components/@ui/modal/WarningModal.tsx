import { Button, Group, Modal, Stack, Text } from '@mantine/core'
import { ReactNode } from 'react'

interface WarningModalProps {
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
  isLoading
}: WarningModalProps) {
  function internalOnYes() {
    onYes()
    onClose()
  }

  return (
    <Modal opened={opened} onClose={onClose} title={title} centered>
      <Modal.Body px={'xs'} py={0}>
        <Stack>
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
      </Modal.Body>
    </Modal>
  )
}

export default WarningModal
