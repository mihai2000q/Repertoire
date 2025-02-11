import { alpha, Avatar, Group, Loader, Menu, Stack, Text, UnstyledButton } from '@mantine/core'
import userPlaceholder from '../../assets/user-placeholder.jpg'
import { IconCaretDownFilled, IconLogout2, IconSettings, IconUser } from '@tabler/icons-react'
import AccountModal from './modal/AccountModal.tsx'
import SettingsModal from './modal/SettingsModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { signOut } from '../../state/slice/authSlice.ts'
import { useGetCurrentUserQuery } from '../../state/api.ts'
import useAuth from '../../hooks/useAuth.ts'
import { useAppDispatch } from '../../state/store.ts'

function TopbarUser() {
  const dispatch = useAppDispatch()

  const { data: user } = useGetCurrentUserQuery(undefined, {
    skip: !useAuth()
  })

  const [openedAccount, { open: openAccount, close: closeAccount }] = useDisclosure(false)
  const [openedSettings, { open: openSettings, close: closeSettings }] = useDisclosure(false)

  function handleSignOut() {
    dispatch(signOut())
  }

  if (!user) return <Loader size={'sm'} />

  return (
    <Menu shadow={'lg'} width={200}>
      <Menu.Target>
        <UnstyledButton
          p={'4px'}
          aria-label={'user'}
          sx={(theme) => ({
            borderRadius: '16px',
            cursor: 'pointer',
            transition: '0.175s all, transform 200ms ease-in-out',
            color: theme.colors.gray[7],
            '&:hover': {
              boxShadow: theme.shadows.sm,
              color: theme.colors.gray[8],
              backgroundColor: alpha(theme.colors.gray[1], 0.7)
            },
            '&:active': {
              transform: 'scale(0.85)'
            }
          })}
        >
          <Group gap={'xxs'}>
            <Avatar src={user.profilePictureUrl ?? userPlaceholder} />
            <IconCaretDownFilled size={12} />
          </Group>
        </UnstyledButton>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Label>
          <Stack gap={0}>
            <Text fw={400} c={'black'}>
              {user.name}
            </Text>
            <Text fz={'xs'} fw={300}>
              {user.email}
            </Text>
          </Stack>
        </Menu.Label>

        <Menu.Divider />

        <Menu.Item leftSection={<IconUser size={14} />} onClick={openAccount}>
          Account
        </Menu.Item>
        <Menu.Item leftSection={<IconSettings size={14} />} onClick={openSettings}>
          Settings
        </Menu.Item>
        <Menu.Item leftSection={<IconLogout2 size={14} />} onClick={handleSignOut}>
          Sign Out
        </Menu.Item>
      </Menu.Dropdown>

      <AccountModal opened={openedAccount} onClose={closeAccount} user={user} />
      <SettingsModal opened={openedSettings} onClose={closeSettings} user={user} />
    </Menu>
  )
}

export default TopbarUser
