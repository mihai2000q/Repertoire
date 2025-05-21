import { ActionIcon, Group } from '@mantine/core'
import { IconChevronLeft, IconChevronRight } from '@tabler/icons-react'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../../../state/store.ts'
import { resetHistoryOnSignIn } from '../../../state/slice/authSlice.ts'
import { useDidUpdate } from '@mantine/hooks'

function TopbarNavigation() {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()

  const historyOnSignIn = useAppSelector((state) => state.auth.historyOnSignIn)

  const disableGoBack = history.state?.idx < 1 + historyOnSignIn.index
  const disableGoForward =
    historyOnSignIn.justSignedIn || history.state?.idx >= history.length - 1

  // when first navigation after sign-in occurs
  useDidUpdate(() => {
    if (historyOnSignIn.justSignedIn) dispatch(resetHistoryOnSignIn())
  }, [history.state?.idx])

  function handleGoBack() {
    navigate(-1)
  }

  function handleGoForward() {
    navigate(1)
  }

  return (
    <Group gap={0} ml={'xs'}>
      <ActionIcon
        aria-label={'back'}
        size={'lg'}
        variant={'grey'}
        radius={'50%'}
        disabled={disableGoBack}
        onClick={handleGoBack}
      >
        <IconChevronLeft size={20} />
      </ActionIcon>

      <ActionIcon
        aria-label={'forward'}
        size={'lg'}
        variant={'grey'}
        radius={'50%'}
        disabled={disableGoForward}
        onClick={handleGoForward}
      >
        <IconChevronRight size={20} />
      </ActionIcon>
    </Group>
  )
}

export default TopbarNavigation
