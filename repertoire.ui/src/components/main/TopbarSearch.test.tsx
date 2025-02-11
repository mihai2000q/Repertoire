import { mantineRender } from '../../test-utils.tsx'
import TopbarSearch from './TopbarSearch.tsx'
import {screen} from "@testing-library/react";

describe('Topbar Search', () => {
  it('should render', () => {
    mantineRender(<TopbarSearch />)

    expect(screen.getByRole('searchbox', { name: /search/i })).toBeInTheDocument()
  })
})
