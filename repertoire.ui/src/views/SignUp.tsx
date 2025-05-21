import { ReactElement } from 'react'
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
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAppDispatch } from '../state/store.ts'
import HttpErrorResponse from '../types/responses/HttpErrorResponse.ts'
import { useForm, zodResolver } from '@mantine/form'
import { signIn } from '../state/slice/authSlice.ts'
import { SignUpForm, signUpValidation } from '../validation/signUpForm.ts'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import { useSignUpMutation } from '../state/api/usersApi.ts'
import { api } from '../state/api.ts'
import { authApi } from '../state/authApi.ts'

function SignUp(): ReactElement {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const location = useLocation()

  useFixedDocumentTitle('Sign Up')

  const [signUpMutation, { error, isLoading }] = useSignUpMutation()
  const signUpError = (error as HttpErrorResponse | undefined)?.data?.error

  const form = useForm({
    mode: 'uncontrolled',
    initialValues: {
      name: '',
      email: '',
      password: ''
    } as SignUpForm,
    validateInputOnBlur: true,
    validateInputOnChange: false,
    clearInputErrorOnChange: true,
    validate: zodResolver(signUpValidation)
  })

  async function signUp({ name, email, password }: SignUpForm): Promise<void> {
    try {
      const token = await signUpMutation({ name, email, password }).unwrap()
      dispatch(signIn(token))
      dispatch(api.util.resetApiState())
      dispatch(authApi.util.resetApiState())
      navigate(location.state?.from?.pathname ?? 'home')
    } catch (e) {
      /*ignored*/
    }
  }

  return (
    <Container h={'100%'}>
      <Center mih={'87vh'} style={{ flexDirection: 'column' }}>
        <Title ta="center" order={2}>
          Create Account
        </Title>
        <Text c="dimmed" size="sm" ta="center" mt={5}>
          Do you have an account already?{' '}
          <Anchor c={'primary.5'} size="sm" component={Link} to={'/sign-in'}>
            Sign in
          </Anchor>
        </Text>

        <Paper withBorder shadow="md" p={30} mt={15}>
          <form onSubmit={form.onSubmit(signUp)}>
            <Stack align={'flex-start'} gap={0} w={200}>
              <Stack w={'100%'}>
                <TextInput
                  label="Name"
                  placeholder="Your name"
                  key={form.key('name')}
                  {...form.getInputProps('name')}
                  {...(signUpError && { error: signUpError })}
                  maxLength={100}
                  disabled={isLoading}
                />
                <TextInput
                  label="Email"
                  placeholder="Your@Email.com"
                  key={form.key('email')}
                  {...form.getInputProps('email')}
                  {...(signUpError && { error: signUpError })}
                  maxLength={256}
                  disabled={isLoading}
                />
                <PasswordInput
                  role={'textbox'}
                  label="Password"
                  placeholder="Your password"
                  key={form.key('password')}
                  {...form.getInputProps('password')}
                  {...(signUpError && { error: signUpError })}
                  disabled={isLoading}
                />
              </Stack>
              <Button type={'submit'} fullWidth mt={'lg'} disabled={isLoading}>
                Sign Up
              </Button>
            </Stack>
          </form>
        </Paper>
      </Center>
    </Container>
  )
}

export default SignUp
