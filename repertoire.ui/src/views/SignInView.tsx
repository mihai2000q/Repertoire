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
import { useSignInMutation } from '../state/api'
import { useAppDispatch } from '../state/store'
import { setToken } from '../state/authSlice'
import HttpErrorResponse from '../types/responses/HttpError.response'
import { useLocation, useNavigate } from 'react-router-dom'

function SignInView(): ReactElement {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const location = useLocation()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const [signInMutation, { error, isLoading }] = useSignInMutation()
  const loginError = (error as HttpErrorResponse | undefined)?.data?.error

  async function signIn(): Promise<void> {
    try {
      const res = await signInMutation({ email, password }).unwrap()

      dispatch(setToken(res.token))
      navigate(location.state?.from?.pathname ?? 'home')
    } catch (e) {
      /*ignored*/
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

        <Paper withBorder shadow="md" p={30} mt={15}>
          <Stack align={'flex-start'} gap={0} w={200}>
            <Stack w={'100%'}>
              <TextInput
                required
                label="Email"
                placeholder="Your email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                error={loginError}
              />
              <PasswordInput
                required
                label="Password"
                placeholder="Your password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                error={loginError}
              />
            </Stack>
            <Anchor component="button" size="sm" mt={6} style={{ alignSelf: 'flex-end' }}>
              Forgot password?
            </Anchor>
            <Button fullWidth mt={'sm'} disabled={isLoading} onClick={signIn}>
              Sign in
            </Button>
          </Stack>
        </Paper>
      </Flex>
    </Container>
  )
}

export default SignInView
