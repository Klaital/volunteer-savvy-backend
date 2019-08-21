require 'test/unit'
require 'json'
require 'net/http'
require 'securerandom'
require 'pp'
require_relative './volunteer_savvy_client.rb'

class TestSiteCrud < Test::Unit::TestCase
    def setup
        @vs_client = VolunteerSavvyClient.new(host: 'localhost', port: ':8080')
        @create_site_data = {
                    slug: "test-create-site",
                    name: "Test Create Site",
                    locale: "en-us",

                    lat: '90.0',
                    lon: '90.0',
                    gplace_id: 'asdfasdf',
                    street: '300 Alamo Plaza',
                    city: 'San Antonio',
                    state: 'TX',
                    zip: '98052',

                }
    end

#    def test_delete_site
#      response = @vs_client.delete_request('/sites/test-create-site/')
#      assert_equal('200', response.code)
      # TODO: actually check if the site was deleted
#    end

    def test_create_site
      response = @vs_client.delete_request('/sites/test-create-site/')
      assert_equal('200', response.code)
      response = @vs_client.post_request('/sites/', @create_site_data.to_json)
      assert_equal('200', response.code)
      # TODO: actually check if the site was created
    end


end
