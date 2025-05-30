import {
  alpha,
  Avatar,
  Center,
  Group,
  Menu,
  Skeleton,
  Stack,
  Text,
  UnstyledButton,
  UnstyledButtonProps
} from '@mantine/core'
import {
  IconCaretDownFilled,
  IconLogout2,
  IconSettings,
  IconUser,
  IconUserFilled
} from '@tabler/icons-react'
import AccountModal from '../modal/AccountModal.tsx'
import SettingsModal from '../modal/SettingsModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { signOut } from '../../../state/slice/authSlice.ts'
import useAuth from '../../../hooks/useAuth.ts'
import { useAppDispatch } from '../../../state/store.ts'
import { useEffect } from 'react'
import { setUserId } from '../../../state/slice/globalSlice.ts'
import { useGetCurrentUserQuery } from '../../../state/api/usersApi.ts'

function TopbarUser({ ...others }: UnstyledButtonProps) {
  const dispatch = useAppDispatch()

  const { data: user } = useGetCurrentUserQuery(undefined, {
    skip: !useAuth()
  })
  useEffect(() => {
    dispatch(setUserId(user?.id))
  }, [dispatch, user])

  const [openedAccount, { open: openAccount, close: closeAccount }] = useDisclosure(false)
  const [openedSettings, { open: openSettings, close: closeSettings }] = useDisclosure(false)

  function handleSignOut() {
    dispatch(signOut())
  }

  if (!user) return <Skeleton w={36} h={36} mx={'xs'} radius={'50%'} style={{ order: 6 }} />

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
          {...others}
        >
          <Group gap={'xxs'}>
            <Avatar src={user.profilePictureUrl} bg={'gray.0'}>
              <Center c={'gray.7'}>
                <IconUserFilled size={20} />
              </Center>
            </Avatar>
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
