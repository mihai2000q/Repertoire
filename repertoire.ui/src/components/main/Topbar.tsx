import { ReactElement } from 'react'
import {
  ActionIcon,
  alpha,
  AppShell,
  Autocomplete,
  Avatar,
  Group,
  Loader,
  Menu,
  Space,
  Stack,
  Text,
  UnstyledButton
} from '@mantine/core'
import userPlaceholder from '../../assets/user-placeholder.jpg'
import {
  IconBellFilled,
  IconCaretDownFilled,
  IconLogout2,
  IconSearch,
  IconSettings,
  IconUser
} from '@tabler/icons-react'
import { useAppDispatch } from '../../state/store.ts'
import { signOut } from '../../state/slice/authSlice.ts'
import { useGetCurrentUserQuery } from '../../state/api.ts'
import useAuth from '../../hooks/useAuth.ts'
import { useDisclosure, useWindowScroll } from '@mantine/hooks'
import AccountModal from './modal/AccountModal.tsx'
import { useNavigate } from 'react-router-dom'
import useIsDesktop from '../../hooks/useIsDesktop.ts'
import CustomIconArrowLeft from '../@ui/icons/CustomIconArrowLeft.tsx'
import CustomIconArrowRight from '../@ui/icons/CustomIconArrowRight.tsx'
import SettingsModal from './modal/SettingsModal.tsx'

function Topbar(): ReactElement {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const isDesktop = useIsDesktop()
  const [scrollPosition] = useWindowScroll()

  const { data: user } = useGetCurrentUserQuery(undefined, {
    skip: !useAuth()
  })

  function handleGoBack() {
    navigate(-1)
  }

  function handleGoForward() {
    navigate(1)
  }

  const [openedAccount, { open: openAccount, close: closeAccount }] = useDisclosure(false)
  const [openedSettings, { open: openSettings, close: closeSettings }] = useDisclosure(false)

  function handleSignOut() {
    dispatch(signOut())
  }

  return (
    <AppShell.Header
      px={'md'}
      withBorder={false}
      top={'unset'}
      style={(theme) => ({
        transition: '0.35s',
        ...(scrollPosition.y !== 0 && { boxShadow: theme.shadows.md })
      })}
    >
      <Group h={'100%'} gap={0}>
        <Autocomplete
          role={'searchbox'}
          aria-label={'topbar-search'}
          placeholder="Search"
          leftSection={<IconSearch size={16} stroke={2} />}
          data={[]}
          fw={500}
          visibleFrom="xs"
          radius={'lg'}
          w={200}
          styles={(theme) => ({
            input: {
              transition: '0.3s',
              backgroundColor: alpha(theme.colors.gray[0], 0.1),
              borderWidth: 0,
              '&:focus, &:hover': {
                boxShadow: theme.shadows.sm,
                backgroundColor: alpha(theme.colors.gray[0], 0.2)
              }
            }
          })}
        />

        {isDesktop && (
          <Group gap={0} ml={'xs'}>
            <ActionIcon
              aria-label={'back-button'}
              size={'lg'}
              variant={'grey'}
              radius={'50%'}
              disabled={window.history.state?.idx < 1}
              onClick={handleGoBack}
            >
              <CustomIconArrowLeft />
            </ActionIcon>

            <ActionIcon
              aria-label={'forward-button'}
              size={'lg'}
              variant={'grey'}
              radius={'50%'}
              disabled={window.history.state?.idx >= window.history.length - 1}
              onClick={handleGoForward}
            >
              <CustomIconArrowRight />
            </ActionIcon>
          </Group>
        )}

        <Space flex={1} />

        <ActionIcon
          variant={'subtle'}
          size={'lg'}
          sx={(theme) => ({
            borderRadius: '50%',
            color: theme.colors.gray[6],
            '&:hover': {
              boxShadow: theme.shadows.sm,
              backgroundColor: theme.colors.primary[0],
              color: theme.colors.primary[6]
            }
          })}
        >
          <IconBellFilled size={18} />
        </ActionIcon>

        {!user ? (
          <Loader size={'sm'} />
        ) : (
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
                <Group gap={4}>
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
        )}
      </Group>
    </AppShell.Header>
  )
}

export default Topbar
