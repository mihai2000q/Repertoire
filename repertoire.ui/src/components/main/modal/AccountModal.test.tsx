import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import { reduxRender, withToastify } from '../../../test-utils.tsx'
import { act, screen } from '@testing-library/react'
import { UpdateAlbumRequest } from '../../../types/requests/AlbumRequests.ts'
import { http, HttpResponse } from 'msw'
import AccountModal from './AccountModal.tsx'
import User from '../../../types/models/User.ts'
import { UpdateUserRequest } from '../../../types/requests/UserRequests.ts'
import dayjs from "dayjs";

describe('Account Modal', () => {
  const user: User = {
    id: '1',
    name: 'User 1',
    email: 'user1@example.com',
    profilePictureUrl: 'some-profile-picture.png',
    createdAt: '2024-11-15T10:30',
    updatedAt: '2024-11-22T11:30'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const userEventDispatcher = userEvent.setup()

    reduxRender(<AccountModal opened={true} onClose={() => {}} user={user} />)

    expect(screen.getByRole('dialog', { name: /account/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /account/i })).toBeInTheDocument()

    expect(screen.getByRole('img', { name: 'profile-picture-preview' })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue(user.name)

    expect(screen.getByRole('textbox', { name: /email/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /email/i })).toBeDisabled()
    expect(screen.getByRole('textbox', { name: /email/i })).toHaveValue(user.email)

    expect(screen.getByText(/created on/i)).toBeInTheDocument()
    expect(screen.getByText(new RegExp(dayjs(user.createdAt).format('DD MMM YYYY')))).toBeInTheDocument()
    expect(screen.getByText(/last modified on/i)).toBeInTheDocument()
    expect(screen.getByText(new RegExp(dayjs(user.updatedAt).format('DD MMM YYYY')))).toBeInTheDocument()

    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toHaveAttribute('data-disabled', 'true')
    await userEventDispatcher.hover(screen.getByRole('button', { name: /save/i }))
    expect(await screen.findByText(/need to make a change/i)).toBeInTheDocument()
  })

  it('should render', async () => {
    const localUser = {
      ...user,
      updatedAt: user.createdAt
    }

    reduxRender(<AccountModal opened={true} onClose={() => {}} user={localUser} />)

    expect(screen.getAllByText(new RegExp(dayjs(localUser.createdAt).format('DD MMM YYYY')))).toHaveLength(1)
    expect(screen.queryByText(/last modified on/i)).not.toBeInTheDocument()
  })

  it('should send only edit request when the profile picture is unchanged', async () => {
    const userEventDispatcher = userEvent.setup()

    const newName = 'New User'
    const onClose = vitest.fn()

    let capturedRequest: UpdateUserRequest
    server.use(
      http.put('/users', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateUserRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<AccountModal opened={true} onClose={onClose} user={user} />))

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    await userEventDispatcher.clear(nameField)
    await userEventDispatcher.type(nameField, newName)
    await userEventDispatcher.click(saveButton)

    expect(await screen.findByText(/account updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      name: newName
    })
  })

  it('should send edit request and save profile picture request when the profile picture is replaced', async () => {
    const userEventDispatcher = userEvent.setup()

    const newImage = new File(['something'], 'profile-picture.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateUserRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/users', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateUserRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/users/pictures', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<AccountModal opened={true} onClose={onClose} user={user} />))

    const saveButton = screen.getByRole('button', { name: /save/i })

    await userEventDispatcher.upload(screen.getByTestId('upload-profile-picture-input'), newImage)
    await userEventDispatcher.click(saveButton)

    expect(await screen.findByText(/account updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      name: user.name
    })
    expect(capturedSaveImageFormData.get('profile_pic')).toBeFormDataImage(newImage)
  })

  it('should send edit request and save profile picture request when the profile picture is first added', async () => {
    const userEventDispatcher = userEvent.setup()

    const newImage = new File(['something'], 'profile-picture.png', { type: 'image/png' })
    const onClose = vitest.fn()

    let capturedRequest: UpdateUserRequest
    let capturedSaveImageFormData: FormData
    server.use(
      http.put('/users', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateUserRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.put('/users/pictures', async (req) => {
        capturedSaveImageFormData = await req.request.formData()
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(
      withToastify(
        <AccountModal
          opened={true}
          onClose={onClose}
          user={{ ...user, profilePictureUrl: undefined }}
        />
      )
    )

    const saveButton = screen.getByRole('button', { name: /save/i })

    await userEventDispatcher.upload(screen.getByTestId('profile-picture-dropzone-input'), newImage)
    await userEventDispatcher.click(saveButton)

    expect(await screen.findByText(/account updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      name: user.name
    })
    expect(capturedSaveImageFormData.get('profile_pic')).toBeFormDataImage(newImage)
  })

  it('should send edit request and delete profile picture request', async () => {
    const userEventDispatcher = userEvent.setup()

    const onClose = vitest.fn()

    let capturedRequest: UpdateAlbumRequest
    server.use(
      http.put('/users', async (req) => {
        capturedRequest = (await req.request.json()) as UpdateAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      }),
      http.delete(`/users/pictures`, () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRender(withToastify(<AccountModal opened={true} onClose={onClose} user={user} />))

    const saveButton = screen.getByRole('button', { name: /save/i })

    await userEventDispatcher.click(screen.getByRole('button', { name: 'remove-profile-picture' }))
    await userEventDispatcher.click(saveButton)

    expect(await screen.findByText(/account updated/i)).toBeInTheDocument()
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    expect(onClose).toHaveBeenCalledOnce()
    expect(capturedRequest).toStrictEqual({
      name: user.name
    })
  })

  it('should disable the save button when no changes are made', async () => {
    const userEventDispatcher = userEvent.setup()

    const newImage = new File(['something'], 'profile-picture.png', { type: 'image/png' })

    reduxRender(<AccountModal opened={true} onClose={() => {}} user={user} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })
    const saveButton = screen.getByRole('button', { name: /save/i })

    // change profile picture
    await userEventDispatcher.upload(screen.getByTestId('upload-profile-picture-input'), newImage)
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset profile picture
    await userEventDispatcher.click(screen.getByRole('button', { name: 'reset-profile-picture' }))
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // change name
    await userEventDispatcher.type(nameField, '1')
    act(() => nameField.blur())
    expect(saveButton).not.toHaveAttribute('data-disabled')

    // reset name
    await userEventDispatcher.clear(nameField)
    await userEventDispatcher.type(nameField, user.name)
    expect(saveButton).toHaveAttribute('data-disabled', 'true')

    // remove profile picture
    await userEventDispatcher.click(screen.getByRole('button', { name: 'remove-profile-picture' }))
    expect(saveButton).not.toHaveAttribute('data-disabled')
  })

  it('should validate the name textbox', async () => {
    const userEventDispatcher = userEvent.setup()

    reduxRender(<AccountModal opened={true} onClose={() => {}} user={user} />)

    const nameField = screen.getByRole('textbox', { name: /name/i })
    await userEventDispatcher.clear(nameField)
    expect(nameField).toBeInvalid()
  })
})
