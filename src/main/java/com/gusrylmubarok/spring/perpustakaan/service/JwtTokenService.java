package com.gusrylmubarok.spring.perpustakaan.service;

import com.gusrylmubarok.spring.perpustakaan.entity.User;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;


@Service
public class JwtTokenService {
    private static final long JWT_TOKEN_VALIDITY = 5 * 60 * 60;
    @Value("${jwt.secret")
    private String secret;

    private Claims getAllClaimsFromToken(String token) {
        return Jwts.parser().setSigningKey(secret).parseClaimsJws(token).getBody();
    }

    private String doGenerateToken(Map<String, Object> claims, String subject) {
        return Jwts.builder().setClaims(claims).setSubject(subject).setIssuedAt((new Date(System.currentTimeMillis())))
                .setExpiration(new Date(System.currentTimeMillis() + JWT_TOKEN_VALIDITY * 1000))
                .signWith(SignatureAlgorithm.HS256, secret).compact();
    }

    public String generateToken(UserDetails userDetails) {
        Map<String, Object> claims = new HashMap();
        return doGenerateToken(claims, userDetails.getUsername());
    }

    public <T> T getClaimsFromToken(String token, Function<Claims, T> claimsTFunction) {
        final Claims claims = getAllClaimsFromToken(token);

        return claimsTFunction.apply(claims);
    }

    public String getUsernameFromToken(String token) {
        return getClaimsFromToken(token, Claims::getSubject);
    }

    public Date getExpirationDateFromToken(String token) {
        return getClaimsFromToken(token, Claims::getExpiration);
    }

    public Boolean isTokenExpired(String token) {
        final Date expiration = getExpirationDateFromToken(token);
        return expiration.before(new Date());
    }

    public Boolean validationToken(String token, User user) {
        return (token.equals(user.getToken()) && !isTokenExpired(token));
    }
}


