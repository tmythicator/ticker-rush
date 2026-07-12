import { describe, it, expect } from 'vitest';
import { ApiError, parseProblemDetails } from './errors';

describe('ApiError', () => {
  it('should initialize correctly with message and status', () => {
    const error = new ApiError('Something went wrong', 400);
    expect(error.message).toBe('Something went wrong');
    expect(error.status).toBe(400);
    expect(error.type).toBeUndefined();
    expect(error.invalidParams).toBeUndefined();
  });

  it('should initialize correctly with RFC 7807 detail parameters', () => {
    const details = {
      type: 'https://tickerrush.com/errors/validation-error',
      title: 'Validation Failed',
      detail: 'Password too short',
      instance: '/api/v1/users',
      invalidParams: [{ name: 'password', reason: 'Must be at least 8 chars' }],
    };

    const error = new ApiError('Password too short', 400, details);
    expect(error.message).toBe('Password too short');
    expect(error.status).toBe(400);
    expect(error.type).toBe(details.type);
    expect(error.title).toBe(details.title);
    expect(error.detail).toBe(details.detail);
    expect(error.instance).toBe(details.instance);
    expect(error.invalidParams).toEqual(details.invalidParams);
  });
});

describe('parseProblemDetails', () => {
  it('should parse RFC 7807 problem details payloads', () => {
    const payload = {
      type: 'https://tickerrush.com/errors/validation-error',
      title: 'Validation Failed',
      detail: 'Password too short',
      instance: '/api/v1/users',
      invalid_params: [{ name: 'password', reason: 'Must be at least 8 chars' }],
    };

    const error = parseProblemDetails(payload, 400);
    expect(error.status).toBe(400);
    expect(error.message).toBe(payload.detail);
    expect(error.type).toBe(payload.type);
    expect(error.title).toBe(payload.title);
    expect(error.detail).toBe(payload.detail);
    expect(error.instance).toBe(payload.instance);
    expect(error.invalidParams).toEqual([{ name: 'password', reason: 'Must be at least 8 chars' }]);
  });

  it('should fallback to legacy error field if RFC 7807 detail is missing', () => {
    const payload = {
      error: 'Legacy error message',
    };

    const error = parseProblemDetails(payload, 400);
    expect(error.message).toBe('Legacy error message');
    expect(error.status).toBe(400);
  });

  it('should fallback to default error message if payload is empty', () => {
    const error = parseProblemDetails(null, 500);
    expect(error.message).toBe('Error: 500');
    expect(error.status).toBe(500);
  });
});
