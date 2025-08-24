package org.example.controller;

import org.example.dto.ResumeUrlResponse;
import org.example.service.S3Service;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.net.URL;
import java.time.Duration;
import java.time.Instant;

@RestController
@RequestMapping("/api/v1")
public class ResumeController {

    private static final String RESUME_OBJECT_KEY = "resume/Resume.pdf";
    private static final Duration DEFAULT_EXPIRATION = Duration.ofMinutes(15);
    private static final Duration MAX_EXPIRATION = Duration.ofHours(24);

    private final S3Service s3Service;

    public ResumeController(S3Service s3Service) {
        this.s3Service = s3Service;
    }

    @GetMapping("/resume")
    public ResponseEntity<ResumeUrlResponse> getResumeUrl(
            @RequestParam(name = "expires_in", required = false) Long expiresInMinutes) {
        
        Duration expiration = DEFAULT_EXPIRATION;
        
        if (expiresInMinutes != null) {
            Duration requestedExpiration = Duration.ofMinutes(expiresInMinutes);
            if (requestedExpiration.compareTo(MAX_EXPIRATION) > 0) {
                expiration = MAX_EXPIRATION;
            } else {
                expiration = requestedExpiration;
            }
        }

        URL signedUrl = s3Service.generateSignedUrl(RESUME_OBJECT_KEY, expiration);
        Instant expiresAt = Instant.now().plus(expiration);

        ResumeUrlResponse response = new ResumeUrlResponse(signedUrl, expiresAt);
        return ResponseEntity.ok(response);
    }
}