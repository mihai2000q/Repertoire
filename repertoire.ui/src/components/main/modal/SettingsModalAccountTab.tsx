import { Button, Divider, Group, Stack, Text } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import DeleteAccountModal from './DeleteAccountModal.tsx'
import User from '../../../types/models/User.ts'

interface SettingsModalAccountTab {
  onCloseSettingsModal: () => void
  user: User
}

function SettingsModalAccountTab({ onCloseSettingsModal, user }: SettingsModalAccountTab) {
  const [openedDelete, { open: openDelete, close: closeDelete }] = useDisclosure(false)

  function handleChangeEmail() {}

  function handleChangePassword() {}

  return (
    <Stack py={'md'} px={'xs'}>
      <Stack gap={0} px={'xs'}>
        <Text fw={600} c={'gray.7'} fz={'xl'}>
          Email
        </Text>
        <Group flex={1} justify={'space-between'}>
          <Text fw={500} c={'dimmed'}>
            Change the email of your account
          </Text>
          <Button variant={'subtle'} onClick={handleChangeEmail}>
            Change Email
          </Button>
        </Group>
      </Stack>

      <Divider />
      <Stack gap={0} px={'xs'}>
        <Text fw={600} c={'gray.7'} fz={'xl'}>
          Password
        </Text>
        <Group flex={1} justify={'space-between'}>
          <Text fw={500} c={'dimmed'}>
            Change the password of your account
          </Text>
          <Button variant={'subtle'} onClick={handleChangePassword}>
            Change Password
          </Button>
        </Group>
      </Stack>

      <Divider />
      <Stack gap={0} px={'xs'}>
        <Text fw={600} c={'gray.7'} fz={'xl'}>
          Deletion
        </Text>
        <Group flex={1} justify={'space-between'}>
          <Text fw={500} c={'dimmed'}>
            This action is a permanent action and cannot be undone
          </Text>
          <Button
            variant={'subtle'}
            sx={(theme) => ({
              color: theme.colors.red[4],
              '&:hover': { color: theme.white, backgroundColor: theme.colors.red[4] }
            })}
            style={{ transition: '0.3s' }}
            onClick={openDelete}
          >
            Delete Account
          </Button>
        </Group>
      </Stack>

      <DeleteAccountModal
        opened={openedDelete}
        onClose={closeDelete}
        onCloseSettingsModal={onCloseSettingsModal}
        user={user}
      />
    </Stack>
  )
}

export default SettingsModalAccountTab
