require 'net/http'
require 'logger'
require 'yaml'
require 'base64'
require 'json'

class VolunteerSavvyClient
  def initialize(opts={})
  puts opts
    @host = opts[:host] || 'localhost'
    @port ||= opts[:port] || ''
    @protocol = opts[:protocol] || (@host == 'localhost') ? 'http' : 'https'
    @logger = opts[:logger] || Logger.new($stdout)
  end

  def construct_uri(endpoint)
    URI("#{@protocol}://#{@host}#{@port}/vs#{endpoint}")
  end
  def get_request(endpoint)
      uri = construct_uri(endpoint)
      # TODO: get login token/cookie

      request = Net::HTTP::Get.new(uri,
                                   'Accept' => 'application/json',
      )
      # @logger.debug "GET #{uri}"
      http_client = Net::HTTP.new(uri.hostname, uri.port)
      http_client.use_ssl = uri.scheme == 'https'
      # http_client.set_debug_output($stderr)
      response = http_client.start do |http|
        http.request(request)
      end

      while response.code == '429'
        sleep(1)
        # logger.debug "Throttled: waiting to retry request #{uri}"
        response = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme == 'https') do |http|
          http.request(request)
        end
      end

      response
    end

    def post_request(endpoint, payload)
      uri = construct_uri(endpoint)
      # TODO: get login token/cookie

      request = Net::HTTP::Post.new(uri,
                                   'Accept' => 'application/json',
                                   'Content-Type' => 'application/json',
      )
      request.body = payload
      # @logger.debug "GET #{uri}"
      http_client = Net::HTTP.new(uri.hostname, uri.port)
      http_client.use_ssl = uri.scheme == 'https'
      # http_client.set_debug_output($stderr)
      response = http_client.start do |http|
        http.request(request)
      end

      while response.code == '429'
        sleep(1)
        # logger.debug "Throttled: waiting to retry request #{uri}"
        response = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme == 'https') do |http|
          http.request(request)
        end
      end

      response
    end
    def put_request(endpoint, payload)
          uri = construct_uri(endpoint)
          # TODO: get login token/cookie

          request = Net::HTTP::Put.new(uri,
                                       'Accept' => 'application/json',
                                       'Content-Type' => 'application/json',
          )
          request.body = payload
          # @logger.debug "GET #{uri}"
          http_client = Net::HTTP.new(uri.hostname, uri.port)
          http_client.use_ssl = uri.scheme == 'https'
          # http_client.set_debug_output($stderr)
          response = http_client.start do |http|
            http.request(request)
          end

          while response.code == '429'
            sleep(1)
            # logger.debug "Throttled: waiting to retry request #{uri}"
            response = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme == 'https') do |http|
              http.request(request)
            end
          end

          response
    end
    def delete_request(endpoint)
      uri = construct_uri(endpoint)
      # TODO: get login token/cookie

      request = Net::HTTP::Delete.new(uri,
      )
      @logger.debug "DELETE  #{uri}"
      http_client = Net::HTTP.new(uri.hostname, uri.port)
      http_client.use_ssl = uri.scheme == 'https'
      # http_client.set_debug_output($stderr)
      response = http_client.start do |http|
        http.request(request)
      end

      while response.code == '429'
        sleep(1)
        # logger.debug "Throttled: waiting to retry request #{uri}"
        response = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme == 'https') do |http|
          http.request(request)
        end
      end

      response
    end
end