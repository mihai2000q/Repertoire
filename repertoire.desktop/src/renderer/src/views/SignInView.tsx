import {
  Anchor,
  Button,
  Container,
  Flex,
  Paper,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Title
} from '@mantine/core'
import { ReactElement, useState } from 'react'
import { useSignInMutation } from '@renderer/state/api'
import { useAppDispatch } from '@renderer/state/store'
import { setToken } from '@renderer/state/authSlice'
import HttpErrorResponse from '@renderer/types/HttpError.response'

function SignInView(): ReactElement {
  const dispatch = useAppDispatch()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const [signIn, { error, isLoading }] = useSignInMutation()
  const loginError = (error as HttpErrorResponse | undefined)?.data

  async function onSubmit(): Promise<void> {
    try {
      const res = await signIn({
        email: email,
        password: password
      }).unwrap()

      dispatch(setToken(res.token))
    } catch (e) {
      window.electron.ipcRenderer.send('log', e)
    }
  }

  return (
    <Container my={40}>
      <Flex direction={'column'} align={'center'} justify={'center'} h={'90%'}>
        <Title ta="center" order={2}>
          Welcome back!
        </Title>
        <Text c="dimmed" size="sm" ta="center" mt={5}>
          Do not have an account yet?{' '}
          <Anchor size="sm" component="button">
            Create account
          </Anchor>
        </Text>

        <Paper withBorder shadow="xl" p={30} mt={15}>
          <form onSubmit={onSubmit}>
            <Stack align={'flex-start'} gap={0} w={200}>
              <Stack w={'100%'}>
                <TextInput
                  required
                  label="Email"
                  placeholder="Your email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  error={loginError?.title}
                />
                <PasswordInput
                  required
                  label="Password"
                  placeholder="Your password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  error={loginError?.title}
                />
              </Stack>
              <Anchor component="button" size="sm" mt={6} style={{ alignSelf: 'flex-end' }}>
                Forgot password?
              </Anchor>
              <Button fullWidth mt={'sm'} type={'submit'} disabled={isLoading}>
                Sign in
              </Button>
            </Stack>
          </form>
        </Paper>
      </Flex>
    </Container>
  )
}

export default SignInView
