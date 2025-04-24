import { reduxRender } from '../../../../../test-utils.tsx'
import { screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { BandMemberRole } from '../../../../../types/models/Artist.ts'
import BandMemberRoleMultiSelect from './BandMemberRoleMultiSelect.tsx'

describe('Band Member Role Multi Select', () => {
  const bandMemberRoles: BandMemberRole[] = [
    {
      id: '1',
      name: 'Guitarist'
    },
    {
      id: '2',
      name: 'Voice'
    },
    {
      id: '3',
      name: 'Drummer'
    }
  ]

  const handlers = [
    http.get('/artists/band-members/roles', async () => {
      return HttpResponse.json(bandMemberRoles)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and change roles', async () => {
    const user = userEvent.setup()

    const newRoles = [bandMemberRoles[0], bandMemberRoles[1]]

    const label = 'band-member-roles'
    const setIds = vitest.fn()

    reduxRender(<BandMemberRoleMultiSelect ids={[]} setIds={setIds} label={label} />)

    const multiSelect = screen.getByRole('textbox', { name: label })
    expect(multiSelect).toHaveValue('')
    expect(multiSelect).toBeDisabled()
    await waitFor(() => expect(multiSelect).not.toBeDisabled())

    await user.click(multiSelect)
    for (const role of bandMemberRoles) {
      expect(await screen.findByRole('option', { name: role.name })).toBeInTheDocument()
    }

    for (const role of newRoles) {
      await user.click(screen.getByRole('option', { name: role.name }))
    }

    expect(setIds).toHaveBeenCalledTimes(newRoles.length)
    newRoles.reduce((a: string[], b) => {
      if (a.length !== 0) expect(setIds).toHaveBeenCalledWith(a) // skip first case
      return [...a, b.id]
    }, [])
  })
})
