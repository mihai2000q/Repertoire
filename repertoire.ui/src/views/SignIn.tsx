import {
  Anchor,
  Button,
  Center,
  Container,
  Paper,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Title
} from '@mantine/core'
import { ReactElement } from 'react'
import { api } from '../state/api'
import { useAppDispatch } from '../state/store'
import { setToken } from '../state/slice/authSlice.ts'
import HttpErrorResponse from '../types/responses/HttpErrorResponse'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useForm, zodResolver } from '@mantine/form'
import { SignInForm, signInValidation } from '../validation/signInForm'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { authApi, useSignInMutation } from '../state/authApi.ts'

function SignIn(): ReactElement {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const location = useLocation()

  useFixedDocumentTitle('Sign In')

  const [signInMutation, { error, isLoading }] = useSignInMutation()
  const signInError = (error as HttpErrorResponse | undefined)?.data?.error

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      email: '',
      password: ''
    } as SignInForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(signInValidation)
  })

  async function signIn({ email, password }: SignInForm): Promise<void> {
    try {
      const token = await signInMutation({ email, password }).unwrap()
      dispatch(setToken(token))
      dispatch(api.util.resetApiState())
      dispatch(authApi.util.resetApiState())
      navigate(location.state?.from?.pathname ?? 'home')
    } catch (e) {
      /*ignored*/
    }
  }

  return (
    <Container h={'100%'}>
      <Center h={'90%'} style={{ flexDirection: 'column' }}>
        <Title ta="center" order={2}>
          Welcome back!
        </Title>
        <Text c="dimmed" size="sm" ta="center" mt={5}>
          Do not have an account yet?{' '}
          <Anchor c={'primary.5'} size="sm" component={Link} to={'/sign-up'}>
            Create account
          </Anchor>
        </Text>

        <Paper withBorder shadow="md" p={30} mt={15}>
          <form onSubmit={form.onSubmit(signIn)}>
            <Stack align={'flex-start'} gap={0} w={200}>
              <Stack w={'100%'}>
                <TextInput
                  label="Email"
                  placeholder="Your email"
                  key={form.key('email')}
                  {...form.getInputProps('email')}
                  {...(signInError && { error: signInError })}
                  maxLength={256}
                  disabled={isLoading}
                />
                <PasswordInput
                  role={'textbox'}
                  label="Password"
                  placeholder="Your password"
                  key={form.key('password')}
                  {...form.getInputProps('password')}
                  {...(signInError && { error: signInError })}
                  disabled={isLoading}
                />
              </Stack>
              <Anchor
                c={'primary.5'}
                component={Link}
                to={'/forgot-pass'}
                size="sm"
                mt={6}
                style={{ alignSelf: 'flex-end' }}
              >
                Forgot password?
              </Anchor>
              <Button type={'submit'} fullWidth mt={'sm'} disabled={isLoading}>
                Sign in
              </Button>
            </Stack>
          </form>
        </Paper>
      </Center>
    </Container>
  )
}

export default SignIn
