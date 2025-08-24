package org.example.dto;

import java.net.URL;
import java.time.Instant;

public record ResumeUrlResponse(URL signedUrl, Instant expiresAt) {
}